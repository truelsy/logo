# logo
*logo*는 *go*언어로 작성된 간단한(Simple) 로깅 라이브러리입니다.

## Installation

    $ go get github.com/truelsy/logo


## QuickStart

```go
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
	logo.CInfof(logo.FG_LIGHT_GREEN, "CInfof(%v)", "logo")

	logo.Warn("Warn")
	logo.Warnf("Warnf(%v)", 789)

	logo.Error("Error")
	logo.Errorf("Error(%v)", 890)
}
```
[![Example Output](example/example.png)]

## Related Projects
- [shiena/ansicolor](https://github.com/shiena/ansicolor)

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License
MIT license
