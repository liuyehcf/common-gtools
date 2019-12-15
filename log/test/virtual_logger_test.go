package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"testing"
	"time"
)

func TestDefaultVirtualLogger(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	assert.AssertTrue(log.GetLogger("notExist") == log.GetLogger("notExist"), "test")

	notExistLogger := log.GetLogger("notExist")

	var content string

	notExistLogger.Trace("you can see this trace log")
	notExistLogger.Debug("you can see this debug log")
	notExistLogger.Info("you can see this info log")
	notExistLogger.Warn("you can see this warn log")
	notExistLogger.Error("you can see this error log")

	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[notExist]-[virtual_logger_test.go:29] --- you can see this info log\n"+
		"[WARN]-[notExist]-[virtual_logger_test.go:30] --- you can see this warn log\n"+
		"[ERROR]-[notExist]-[virtual_logger_test.go:31] --- you can see this error log\n", content)

	time.Sleep(time.Second)
}

func TestChangeTargetLogger(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	logger := log.GetLogger(log.Root)

	var content string

	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertFalse(logger.IsDebugEnabled(), "test")
	assert.AssertTrue(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")
	logger.Trace("you cannot see this trace log")
	logger.Debug("you cannot see this debug log")
	logger.Info("you can see this info log")
	logger.Warn("you can see this warn log")
	logger.Error("you can see this error log")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[ROOT]-[virtual_logger_test.go:63] --- you can see this info log\n"+
		"[WARN]-[ROOT]-[virtual_logger_test.go:64] --- you can see this warn log\n"+
		"[ERROR]-[ROOT]-[virtual_logger_test.go:65] --- you can see this error log\n", content)

	newLogger := log.NewLogger(log.Root, log.TraceLevel, false, []log.Appender{writerAppender})

	logger.Trace("you can see this trace log")
	newLogger.Trace("you can see this trace log")
	logger.Debug("you can see this debug log")
	newLogger.Debug("you can see this debug log")
	logger.Info("you can see this info log")
	newLogger.Info("you can see this info log")
	logger.Warn("you can see this warn log")
	newLogger.Warn("you can see this warn log")
	logger.Error("you can see this error log")
	newLogger.Error("you can see this error log")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[WARN]-[ROOT]-[logger.go:209] --- logger 'ROOT' is replaced\n"+
		"[TRACE]-[ROOT]-[virtual_logger_test.go:74] --- you can see this trace log\n"+
		"[TRACE]-[ROOT]-[virtual_logger_test.go:75] --- you can see this trace log\n"+
		"[DEBUG]-[ROOT]-[virtual_logger_test.go:76] --- you can see this debug log\n"+
		"[DEBUG]-[ROOT]-[virtual_logger_test.go:77] --- you can see this debug log\n"+
		"[INFO]-[ROOT]-[virtual_logger_test.go:78] --- you can see this info log\n"+
		"[INFO]-[ROOT]-[virtual_logger_test.go:79] --- you can see this info log\n"+
		"[WARN]-[ROOT]-[virtual_logger_test.go:80] --- you can see this warn log\n"+
		"[WARN]-[ROOT]-[virtual_logger_test.go:81] --- you can see this warn log\n"+
		"[ERROR]-[ROOT]-[virtual_logger_test.go:82] --- you can see this error log\n"+
		"[ERROR]-[ROOT]-[virtual_logger_test.go:83] --- you can see this error log\n", content)

	assert.AssertTrue(logger == newLogger, "test")
}

func TestLoggerLevel(t *testing.T) {
	logger := log.GetLogger("notExist")

	log.NewLogger(log.Root, log.TraceLevel, false, nil)
	assert.AssertTrue(logger.IsTraceEnabled(), "test")
	assert.AssertTrue(logger.IsDebugEnabled(), "test")
	assert.AssertTrue(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")

	log.NewLogger(log.Root, log.DebugLevel, false, nil)
	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertTrue(logger.IsDebugEnabled(), "test")
	assert.AssertTrue(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")

	log.NewLogger(log.Root, log.InfoLevel, false, nil)
	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertFalse(logger.IsDebugEnabled(), "test")
	assert.AssertTrue(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")

	log.NewLogger(log.Root, log.WarnLevel, false, nil)
	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertFalse(logger.IsDebugEnabled(), "test")
	assert.AssertFalse(logger.IsInfoEnabled(), "test")
	assert.AssertTrue(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")

	log.NewLogger(log.Root, log.ErrorLevel, false, nil)
	assert.AssertFalse(logger.IsTraceEnabled(), "test")
	assert.AssertFalse(logger.IsDebugEnabled(), "test")
	assert.AssertFalse(logger.IsInfoEnabled(), "test")
	assert.AssertFalse(logger.IsWarnEnabled(), "test")
	assert.AssertTrue(logger.IsErrorEnabled(), "test")
}
