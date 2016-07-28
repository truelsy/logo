// logo.go
package logo

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shiena/ansicolor"
)

type Environment struct {
	LogLevel       int
	LogPath        string
	LogKeepTime    time.Duration
	RotateFileSize int
	WriteConsole   bool
}

func (env *Environment) init() {
	if LEVEL_DEBUG > env.LogLevel || env.LogLevel > LEVEL_FATAL {
		env.LogLevel = LEVEL_DEBUG
	}

	if "" == env.LogPath {
		env.LogPath = "."
	} else {
		os.MkdirAll(env.LogPath, 0755)
	}

	//if 0 == env.LogKeepTime {
	//	env.LogKeepTime = time.Hour * 24 * 15 // 15DAY
	//}

	if 0 == env.RotateFileSize {
		env.RotateFileSize = (1024 << 10) * 10 // 10MB
	}
}

type Logger struct {
	rotateFileNum int
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	curFile       *os.File
	filePrefix    string
	env           *Environment
	sync.Mutex
}

var reFileName *regexp.Regexp
var reFuncName *regexp.Regexp
var loggerMap map[int]*Logger = map[int]*Logger{}

func Init(env *Environment) {
	env.init()

	for level := LEVEL_DEBUG; level <= LEVEL_FATAL; level++ {
		logger := new(Logger)
		logger.env = env
		switch level {
		case LEVEL_DEBUG:
			logger.filePrefix = "DEBUG"
		case LEVEL_INFO:
			logger.filePrefix = "INFO"
		case LEVEL_WARN:
			logger.filePrefix = "WARN"
		case LEVEL_ERROR:
			logger.filePrefix = "ERROR"
		case LEVEL_FATAL:
			logger.filePrefix = "FATAL"
		default:
			logger.filePrefix = "UNKNOWN"
		}

		if env.WriteConsole {
			// create console logger
			logger.consoleLogger = log.New(ansicolor.NewAnsiColorWriter(os.Stdout), "", PrintLogFormat)
		}

		loggerMap[level] = logger
	}

	reFileName = regexp.MustCompile("[a-zA-Z]+_[0-9]{8}_[0-9]{6}")
	reFuncName = regexp.MustCompile(`^.*\.(.*)$`)
}

func (logger *Logger) removeOldFile() {
	if 0 >= logger.env.LogKeepTime {
		return
	}

	fileList, err := filepath.Glob(path.Join(logger.env.LogPath, logger.filePrefix+"_*.log"))
	if nil != err {
		return
	}

	now := time.Now()
	for _, fileName := range fileList {
		fileInfo, err := os.Stat(fileName)
		if nil != err {
			continue
		}
		if logger.env.LogKeepTime <= now.Sub(fileInfo.ModTime()) {
			os.Remove(fileName)
		}
	}
}

func (logger *Logger) getNewFileName() string {
	now := time.Now()

	fileName := fmt.Sprintf("%s_%d%02d%02d_%02d%02d%02d",
		logger.filePrefix,
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second())

	return fileName
}

func (logger *Logger) getCurrentFileName() string {
	if nil == logger.curFile {
		return ""
	}

	// remove extention
	_, fileName := filepath.Split(logger.curFile.Name())
	extension := filepath.Ext(fileName)
	fileName = fileName[0 : len(fileName)-len(extension)]

	m := reFileName.FindAllString(fileName, -1)
	if nil == m {
		return ""
	}
	return m[0]
}

func (logger *Logger) createNewFile() {
	newFileName := logger.getNewFileName()
	curFileName := logger.getCurrentFileName()

	if 0 == strings.Compare(newFileName, curFileName) {
		logger.rotateFileNum++
		newFileName = fmt.Sprintf("%s_%d", newFileName, logger.rotateFileNum)
	} else {
		logger.rotateFileNum = 0
	}

	f, err := os.OpenFile(path.Join(logger.env.LogPath, newFileName+".log"), FileOpenFlags, 0666)
	if err != nil {
		return
	}

	// open logger close
	logger.Close()

	//	if logger.env.WriteConsole {
	//		logger.fileLogger = log.New(io.MultiWriter(f, os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds)
	//	} else {
	//		logger.fileLogger = log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	//	}

	logger.fileLogger = log.New(f, "", PrintLogFormat)
	logger.curFile = f
}

func (logger *Logger) isRotateFile() bool {
	fileInfo, err := logger.curFile.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "isRotateFile: file.Stat Failed.")
		return false
	}

	if fileInfo.Size() >= int64(logger.env.RotateFileSize) {
		return true
	}
	return false
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.curFile != nil {
		logger.curFile.Close()
	}

	logger.fileLogger = nil
	logger.curFile = nil
}

func (logger *Logger) doPrint(level int, printLevel string, color int, format string, args ...interface{}) {
	logger.Lock()
	defer logger.Unlock()

	if level < logger.env.LogLevel {
		return
	}

	logger.removeOldFile()

	if logger.fileLogger == nil {
		logger.createNewFile()
	}

	if logger.isRotateFile() {
		logger.createNewFile()
	}

	var msg string

	if "" == format {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	// no log message
	if "" == msg {
		return
	}

	pc, fileName, fileLine, _ := runtime.Caller(3)
	funcName := reFuncName.ReplaceAllString(runtime.FuncForPC(pc).Name(), "$1")

	fileInfo := fmt.Sprintf(PrintFileInfoFormat, filepath.Base(fileName), fileLine, funcName)

	// print to console
	if nil != logger.consoleLogger {
		if COLOR_NONE < color {
			msgColor := getColor(color) + fileInfo + msg + "\033[0m"
			logger.consoleLogger.Print(printLevel, msgColor)
		} else {
			logger.consoleLogger.Print(printLevel, fileInfo+msg)
		}
		//os.Stdout.Sync()
	}

	// print to writer
	logger.fileLogger.Print(printLevel, fileInfo+msg)
	//logger.curFile.Sync()

	if level == LEVEL_FATAL {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, color int, args ...interface{}) {
	logger.doPrint(LEVEL_DEBUG, PrintDebug, color, format, args...)
}

func (logger *Logger) Info(format string, color int, args ...interface{}) {
	logger.doPrint(LEVEL_INFO, PrintInfo, color, format, args...)
}

func (logger *Logger) Warn(format string, color int, args ...interface{}) {
	logger.doPrint(LEVEL_WARN, PrintWarn, color, format, args...)
}

func (logger *Logger) Error(format string, color int, args ...interface{}) {
	logger.doPrint(LEVEL_ERROR, PrintError, color, format, args...)
}

func (logger *Logger) Fatal(format string, color int, args ...interface{}) {
	logger.doPrint(LEVEL_FATAL, PrintFatal, color, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if logger := getLogger(LEVEL_DEBUG); nil != logger {
		logger.Debug(format, COLOR_NONE, args...)
	}
}

func Debug(args ...interface{}) {
	if logger := getLogger(LEVEL_DEBUG); nil != logger {
		logger.Debug("", COLOR_NONE, args...)
	}
}

func CDebugf(color int, format string, args ...interface{}) {
	if logger := getLogger(LEVEL_DEBUG); nil != logger {
		logger.Debug(format, color, args...)
	}
}

func CDebug(color int, args ...interface{}) {
	if logger := getLogger(LEVEL_DEBUG); nil != logger {
		logger.Debug("", color, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger := getLogger(LEVEL_INFO); nil != logger {
		logger.Info(format, COLOR_NONE, args...)
	}
}

func Info(args ...interface{}) {
	if logger := getLogger(LEVEL_INFO); nil != logger {
		logger.Info("", COLOR_NONE, args...)
	}
}

func CInfof(color int, format string, args ...interface{}) {
	if logger := getLogger(LEVEL_INFO); nil != logger {
		logger.Info(format, color, args...)
	}
}

func CInfo(color int, args ...interface{}) {
	if logger := getLogger(LEVEL_INFO); nil != logger {
		logger.Info("", color, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger := getLogger(LEVEL_WARN); nil != logger {
		logger.Warn(format, FG_CYAN, args...)
	}
}

func Warn(args ...interface{}) {
	if logger := getLogger(LEVEL_WARN); nil != logger {
		logger.Warn("", FG_CYAN, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger := getLogger(LEVEL_ERROR); nil != logger {
		logger.Error(format, FG_RED, args...)
	}
}

func Error(args ...interface{}) {
	if logger := getLogger(LEVEL_ERROR); nil != logger {
		logger.Error("", FG_RED, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger := getLogger(LEVEL_FATAL); nil != logger {
		logger.Fatal(format, FG_RED, args...)
	}
}

func Fatal(args ...interface{}) {
	if logger := getLogger(LEVEL_FATAL); nil != logger {
		logger.Fatal("", FG_RED, args...)
	}
}

func Close() {
	for _, logger := range loggerMap {
		logger.Close()
	}
}

func getLogger(level int) *Logger {
	logger, ok := loggerMap[level]
	if ok {
		return logger
	}
	return nil
}

func getColor(color int) string {
	switch color {
	case FG_BLACK:
		return "\033[30m"
	case FG_RED:
		return "\033[31m"
	case FG_GREEN:
		return "\033[32m"
	case FG_YELLOW:
		return "\033[33m"
	case FG_BLUE:
		return "\033[34m"
	case FG_MAGENTA:
		return "\033[35m"
	case FG_CYAN:
		return "\033[36m"
	case FG_WHITE:
		return "\033[37m"
	case FG_DEFAULT:
		return "\033[39m"

	case FG_LIGHT_GRAY:
		return "\033[90m"
	case FG_LIGHT_RED:
		return "\033[91m"
	case FG_LIGHT_GREEN:
		return "\033[92m"
	case FG_LIGHT_YELLOW:
		return "\033[93m"
	case FG_LIGHT_BLUE:
		return "\033[94m"
	case FG_LIGHT_MAGENTA:
		return "\033[95m"
	case FG_LIGHT_CYAN:
		return "\033[96m"
	case FG_LIGHT_WHITE:
		return "\033[97m"

	case BG_BLACK:
		return "\033[40m"
	case BG_RED:
		return "\033[41m"
	case BG_GREEN:
		return "\033[42m"
	case BG_YELLOW:
		return "\033[43m"
	case BG_BLUE:
		return "\033[44m"
	case BG_MAGENTA:
		return "\033[45m"
	case BG_CYAN:
		return "\033[46m"
	case BG_WHITE:
		return "\033[47m"
	case BG_DEFAULT:
		return "\033[49m"

	case BG_LIGHT_GRAY:
		return "\033[100m"
	case BG_LIGHT_RED:
		return "\033[101m"
	case BG_LIGHT_GREEN:
		return "\033[102m"
	case BG_LIGHT_YELLOW:
		return "\033[103m"
	case BG_LIGHT_BLUE:
		return "\033[104m"
	case BG_LIGHT_MAGENTA:
		return "\033[105m"
	case BG_LIGHT_CYAN:
		return "\033[106m"
	case BG_LIGHT_WHITE:
		return "\033[107m"

		//	case FG_FG_INTENSITY_ON:
		//		return "\033[1m"
		//	case FG_INTENSITY_OFF:
		//		return "\033[21m"
		//	case FG_UNDERLINE_ON:
		//		return "\033[4m"
		//	case FG_UNDERLINE_OFF:
		//		return "\033[24m"
		//	case FG_BLINK_ON:
		//		return "\033[5m"
		//	case FG_BLINK_OFF:
		//		return "\033[25m"
	default:
		return ""
	}
}
