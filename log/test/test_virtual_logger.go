package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/log"
	"os"
	"time"
)

func main() {
	assert.AssertTrue(log.GetLogger("something") == log.GetLogger("something"), "test")

	logger := log.GetLogger(log.Root)

	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertFalse(logger.IsDebugEnabled(), "test")
	assert.AssertTrue(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")
	logger.Trace("you can see this log1")
	logger.Debug("you can see this log1")
	logger.Info("you can see this log1")
	logger.Warn("you can see this log1")
	logger.Error("you can see this log1")

	layout := "%d{2006-01-02 15:04:05.999} [%-5p] [%L] %m%n"
	stdoutAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  layout,
		Filters: nil,
		Writer:  os.Stdout,
	})

	newLogger := log.NewLogger(log.Root, log.TraceLevel, false, []log.Appender{stdoutAppender})

	logger.Trace("you can see this trace log2")
	newLogger.Trace("you can see this trace log2")
	logger.Debug("you can see this debug log2")
	newLogger.Debug("you can see this debug log2")
	logger.Info("you can see this info log2")
	newLogger.Info("you can see this info log2")
	logger.Warn("you can see this warn log2")
	newLogger.Warn("you can see this warn log2")
	logger.Error("you can see this error log2")
	newLogger.Error("you can see this error log2")

	assert.AssertTrue(logger == newLogger, "test")

	notExistLogger := log.GetLogger("notExist")

	notExistLogger.Trace("you can see this trace log3")
	notExistLogger.Debug("you can see this debug log3")
	notExistLogger.Info("you can see this info log3")
	notExistLogger.Warn("you can see this warn log3")
	notExistLogger.Error("you can see this error log3")

	time.Sleep(time.Second)
}
