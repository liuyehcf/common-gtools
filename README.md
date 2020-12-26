# Buffer

```go
package main

import (
	"fmt"
	buf "github.com/liuyehcf/common-gtools/buffer"
)

func main() {
	buffer := buf.NewRecycleByteBuffer(10)

	buffer.Write([]byte{1, 2, 3, 4, 5})

	buffer.Mark()

	fmt.Printf("after write, readableBytes=%d\n", buffer.ReadableBytes())
	bytes := make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())

	buffer.Recover()
	fmt.Printf("after recover, readableBytes=%d\n", buffer.ReadableBytes())
	bytes = make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())

	buffer.Write([]byte{6, 7, 8, 9, 10})
	fmt.Printf("after write, readableBytes=%d\n", buffer.ReadableBytes())
	bytes = make([]byte, 5)
	buffer.Read(bytes)
	fmt.Println(bytes)
	fmt.Printf("after read, readableBytes=%d\n", buffer.ReadableBytes())
}
```

# Log

Go version of logback, which is the main log framework of Java

| conversion | description |
|:--|:--|
| `c`/`lo`/`logger` | logger name<br>support left and right alignment and width setting  |
| `d{format}`/`date{format}` | date in specified go's time format, like `2006-01-02 15:04:05.999`<br>support left and right alignment and width setting |
| `L`/`line` | simple source file name and line num, like `main.go:34`<br>support left and right alignment and width setting |
| `m`/`msg`/`message` | log message<br>support left and right alignment and width setting |
| `n` | new line |
| `p`/`le`/`level` | log level, including `TRACE`、`DEBUG`、`INFO`、`WARN`、`ERROR`<br>support left and right alignment and width setting |

```go
package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

var logger = log.GetLogger("testLogger")

func main() {
	logger.Trace("current time is {}", time.Now())
	logger.Debug("current time is {}", time.Now())
	logger.Info("current time is {}", time.Now())
	logger.Warn("current time is {}", time.Now())
	logger.Error("current time is {}", time.Now())
	time.Sleep(time.Second)
}

func init() {
	layout := "%-24d{2006-01-02 15:04:05.999} [%-10c] [%-5p] --- [%L] %m%n"
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	stdoutAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{infoLevelFilter},
		Writer:  os.Stdout,
	})
	stderrAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{errorLevelFilter},
		Writer:  os.Stderr,
	})

	commonFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{infoLevelFilter},
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs",
			FileName:        "common",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})

	errorFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{errorLevelFilter},
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs",
			FileName:        "error",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{commonFileAppender, errorFileAppender, stdoutAppender, stderrAppender})
}
```