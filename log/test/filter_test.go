package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"testing"
	"time"
)

func TestNoFilter(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[ROOT]-[filter_test.go:23] --- you can see this once\n", content)
}

func TestTwoSameFilters(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	infoLevelFilter1 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	infoLevelFilter2 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: []log.Filter{infoLevelFilter1, infoLevelFilter2},
		Writer:  writer,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you can see this twice", time.Now())
	time.Sleep(time.Millisecond * 10)
	logger.Error("you can see this twice", time.Now())
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[ROOT]-[filter_test.go:47] --- you can see this twice\n"+
		"[ERROR]-[ROOT]-[filter_test.go:49] --- you can see this twice\n", content)
}

func TestTwoDifferentFilters(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: []log.Filter{infoLevelFilter, errorLevelFilter},
		Writer:  writer,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you cannot see this once", time.Now())
	logger.Error("you can see this once", time.Now())
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[ERROR]-[ROOT]-[filter_test.go:75] --- you can see this once\n", content)
}

func TestNilFilter(t *testing.T) {
	var nilLevelFilter *log.LevelFilter = nil
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: []log.Filter{nil, nilLevelFilter},
		Writer:  writer,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[ROOT]-[filter_test.go:94] --- you can see this once\n", content)
}
