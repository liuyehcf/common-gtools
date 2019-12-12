package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

func main() {
	leftAlign := "%-30d{2006-01-02 15:04:05.999} [%-10p] %-1m%n"
	rightAlign := "%30d{2006-01-02 15:04:05.999} [%10p] %1m%n"
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  leftAlign,
		Filters: []log.Filter{infoLevelFilter},
		Writer:  os.Stdout,
	})
	stderrAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  rightAlign,
		Filters: []log.Filter{errorLevelFilter},
		Writer:  os.Stderr,
	})

	commonFileAppender := log.NewFileAppender(&log.AppenderConfig{
		Layout:  leftAlign,
		Filters: []log.Filter{infoLevelFilter},
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs",
			FileName:        "common",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})

	errorFileAppender := log.NewFileAppender(&log.AppenderConfig{
		Layout:  rightAlign,
		Filters: []log.Filter{errorLevelFilter},
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs",
			FileName:        "error",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})

	logger := log.DefaultLogger{
		Level:     log.InfoLevel,
		Appenders: []log.Appender{commonFileAppender, errorFileAppender, stdoutAppender, stderrAppender},
	}

	go func() {
		for {
			logger.Trace("现在的时间是 {}", time.Now())
			logger.Debug("现在的时间是 {}", time.Now())
			logger.Info("现在的时间是 {}", time.Now())
			logger.Warn("现在的时间是 {}", time.Now())
			logger.Error("现在的时间是 {}", time.Now())

			time.Sleep(1 * time.Second)
		}
	}()

	<-make(chan interface{}, 0)
}
