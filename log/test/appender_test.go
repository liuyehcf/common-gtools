package main

import (
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"github.com/liuyehcf/common-gtools/utils"
	"testing"
	"time"
)

func TestNilAppender(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	commonFileAppender, _ := log.NewFileAppender(&log.AppenderConfig{
		Layout: "[%p]-[%c]-[%L] --- %m%n",
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/a/b/c",
			FileName:        "common",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})
	utils.AssertTrue(utils.IsNil(commonFileAppender), "test")

	logger := log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender, commonFileAppender, nil})

	var content string

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "[INFO]-[ROOT]-[appender_test.go:35] --- you can see this once\n", content)
}
