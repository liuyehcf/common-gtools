package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/log"
	"testing"
)

func Test(t *testing.T) {
	leftAlign := "%-30d{2006-01-02 15:04:05.999} [%-10c] [%-10p] --- [%-20L] %-1m%n"
	commonFileAppender, err := log.NewFileAppender(&log.AppenderConfig{
		Layout: leftAlign,
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs/",
			FileName:        "common",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})
	assert.AssertNotNil(commonFileAppender, "test")
	assert.AssertNil(err, "test")

	commonFileAppender, err = log.NewFileAppender(&log.AppenderConfig{
		Layout: leftAlign,
		FileRollingPolicy: &log.RollingPolicy{
			Directory:       "/tmp/gtools/logs///",
			FileName:        "common",
			TimeGranularity: log.TimeGranularityHour,
			MaxHistory:      10,
			MaxFileSize:     1024 * 1024 * 1024,
		},
	})
	assert.AssertNotNil(commonFileAppender, "test")
	assert.AssertNil(err, "test")

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{commonFileAppender})

	logger := log.GetLogger("test")
	logger.Info("hello world")
}
