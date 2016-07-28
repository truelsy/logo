// const.go
package logo

import (
	"log"
	"os"
)

// console color
const (
	COLOR_NONE = iota //0

	FG_BLACK
	FG_RED
	FG_GREEN
	FG_YELLOW
	FG_BLUE
	FG_MAGENTA
	FG_CYAN
	FG_WHITE
	FG_DEFAULT

	FG_LIGHT_GRAY
	FG_LIGHT_RED
	FG_LIGHT_GREEN
	FG_LIGHT_YELLOW
	FG_LIGHT_BLUE
	FG_LIGHT_MAGENTA
	FG_LIGHT_CYAN
	FG_LIGHT_WHITE

	BG_BLACK
	BG_RED
	BG_GREEN
	BG_YELLOW
	BG_BLUE
	BG_MAGENTA
	BG_CYAN
	BG_WHITE
	BG_DEFAULT

	BG_LIGHT_GRAY
	BG_LIGHT_RED
	BG_LIGHT_GREEN
	BG_LIGHT_YELLOW
	BG_LIGHT_BLUE
	BG_LIGHT_MAGENTA
	BG_LIGHT_CYAN
	BG_LIGHT_WHITE
)

// levels
const (
	LEVEL_DEBUG = iota // 0
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

const (
	PrintDebug = "[DEBUG ] - "
	PrintInfo  = "[INFO  ] - "
	PrintWarn  = "[WARN  ] - "
	PrintError = "[ERROR ] - "
	PrintFatal = "[FATAL ] - "

	//PrintTimeFormat = "2006-01-02 15:04:05.000000 "
	PrintLogFormat      = log.Ldate | log.Ltime /*| log.Lmicroseconds*/
	PrintFileInfoFormat = "[%-14s:%-4d] %s : "
	FileOpenFlags       = os.O_WRONLY | os.O_CREATE | os.O_EXCL /*| os.O_SYNC*/
)
