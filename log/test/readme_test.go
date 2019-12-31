package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"testing"
	"time"
)

var logger = log.GetLogger("testLogger")

func TestReadMe(t *testing.T) {
	logger.Trace("current time is {}", time.Now())
	logger.Debug("current time is {}", time.Now())
	logger.Info("current time is {}", time.Now())
	logger.Warn("current time is {}", time.Now())
	logger.Error("current time is {}", time.Now())
	time.Sleep(1)
}

func init() {
	leftAlign := "%-30d{2006-01-02 15:04:05.999} [%-10c] [%-10p] --- [%-20L] %-1m%n"
	rightAlign := "%30d{2006-01-02 15:04:05.999} [%10c] [%10p] --- [%20L] %1m%n"
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	stdoutAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  leftAlign,
		Filters: []log.Filter{infoLevelFilter},
		Writer:  os.Stdout,
	})
	stderrAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  rightAlign,
		Filters: []log.Filter{errorLevelFilter},
		Writer:  os.Stderr,
	})

	commonFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
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

	errorFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
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

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{commonFileAppender, errorFileAppender, stdoutAppender, stderrAppender})
}
