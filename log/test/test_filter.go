package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"time"
)

func main() {
	testNoFilter()
	twoSameFilters()
	twoDifferentFilters()
}

func testNoFilter() {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[ROOT]-[test_filter.go:28] --- you can see this once\n", content)
}

func twoSameFilters() {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	infoLevelFilter1 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	infoLevelFilter2 := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
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
	assert.AssertTrue(content == "[INFO]-[ROOT]-[test_filter.go:52] --- you can see this twice\n"+
		"[ERROR]-[ROOT]-[test_filter.go:54] --- you can see this twice\n", content)
}

func twoDifferentFilters() {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	infoLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.InfoLevel,
	}
	errorLevelFilter := &log.LevelFilter{
		LogLevelThreshold: log.ErrorLevel,
	}
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
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
	assert.AssertTrue(content == "[ERROR]-[ROOT]-[test_filter.go:80] --- you can see this once\n", content)
}
