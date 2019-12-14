package main

import (
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

func main() {
	testNoFilter()
	twoSameFilters()
	twoDifferentFilters()
}

func testNoFilter() {
	layout := "%d{2006-01-02 15:04:05.999} [%p] %m%n"
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: nil,
		Writer:  os.Stdout,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{stdoutAppender})

	logger.Info("you can see this1 once, now={}", time.Now())

	time.Sleep(time.Second)
}

func twoSameFilters() {
	layout := "%d{2006-01-02 15:04:05.999} [%p] %m%n"
	infoLevelFilter1 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	infoLevelFilter2 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: []log.Filter{infoLevelFilter1, infoLevelFilter2},
		Writer:  os.Stdout,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{stdoutAppender})

	logger.Info("you can see this2 twice, now={}", time.Now())
	logger.Error("you can see this2 twice, now={}", time.Now())

	time.Sleep(time.Second)
}

func twoDifferentFilters() {
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

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{stdoutAppender})

	logger.Info("you cannot see this3 once, now={}", time.Now())
	logger.Error("you can see this3 once, now={}", time.Now())

	time.Sleep(time.Second)
}
