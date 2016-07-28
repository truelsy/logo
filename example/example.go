// example.go
package main

import (
	"time"

	"github.com/truelsy/logo"
)

func main() {
	logo.Init(&logo.Environment{
		LogLevel:       logo.LEVEL_DEBUG,
		LogPath:        "/tmp/logo/",
		RotateFileSize: 1024 * 10, // 10K
		LogKeepTime:    time.Hour * 6,
		WriteConsole:   true,
	})
	defer logo.Close()

	logo.Debug("Debug")
	logo.Debugf("Debugf(%v)", 123)
	logo.CDebug(logo.FG_GREEN, "CDebug")
	logo.CDebugf(logo.FG_LIGHT_GREEN, "CDebugf(%v)", "hello")

	logo.Info("Info")
	logo.Infof("Infof(%v)", 456)
	logo.CInfo(logo.FG_YELLOW, "CInfo")
	logo.CInfof(logo.FG_LIGHT_YELLOW, "CInfof(%v)", "logo")

	logo.Warn("Warn")
	logo.Warnf("Warnf(%v)", 789)

	logo.Error("Error")
	logo.Errorf("Error(%v)", 890)
}
