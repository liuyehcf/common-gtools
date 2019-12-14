package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

func main() {
	testNoFilter()

	twoFilters()
}

func testNoFilter() {
	layout := "%d{2006-01-02 15:04:05.999} [%p] %m%n"
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: nil,
		Writer:  os.Stdout,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{stdoutAppender})

	logger.Info("test now={}", time.Now())

	time.Sleep(time.Second * 2)
}

func twoFilters() {
	layout := "%d{2006-01-02 15:04:05.999} [%p] %m%n"
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{infoLevelFilter, errorLevelFilter},
		Writer:  os.Stdout,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, []log.Appender{stdoutAppender})

	logger.Info("test now={}", time.Now())
	logger.Error("test now={}", time.Now())

	time.Sleep(time.Second * 2)
}
