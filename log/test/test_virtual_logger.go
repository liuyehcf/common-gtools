package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"time"
)

func main() {
	assert.AssertTrue(log.GetLogger("something") == log.GetLogger("something"), "test")

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
	assert.AssertTrue(content == "[INFO]-[ROOT]-[test_virtual_logger.go:33] --- you can see this info log\n"+
		"[WARN]-[ROOT]-[test_virtual_logger.go:34] --- you can see this warn log\n"+
		"[ERROR]-[ROOT]-[test_virtual_logger.go:35] --- you can see this error log\n", content)

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
	assert.AssertTrue(content == "[WARN]-[ROOT]-[logger.go:186] --- logger 'ROOT' is replaced\n"+
		"[TRACE]-[ROOT]-[test_virtual_logger.go:44] --- you can see this trace log\n"+
		"[TRACE]-[ROOT]-[test_virtual_logger.go:45] --- you can see this trace log\n"+
		"[DEBUG]-[ROOT]-[test_virtual_logger.go:46] --- you can see this debug log\n"+
		"[DEBUG]-[ROOT]-[test_virtual_logger.go:47] --- you can see this debug log\n"+
		"[INFO]-[ROOT]-[test_virtual_logger.go:48] --- you can see this info log\n"+
		"[INFO]-[ROOT]-[test_virtual_logger.go:49] --- you can see this info log\n"+
		"[WARN]-[ROOT]-[test_virtual_logger.go:50] --- you can see this warn log\n"+
		"[WARN]-[ROOT]-[test_virtual_logger.go:51] --- you can see this warn log\n"+
		"[ERROR]-[ROOT]-[test_virtual_logger.go:52] --- you can see this error log\n"+
		"[ERROR]-[ROOT]-[test_virtual_logger.go:53] --- you can see this error log\n", content)

	assert.AssertTrue(logger == newLogger, "test")

	notExistLogger := log.GetLogger("notExist")

	notExistLogger.Trace("you can see this trace log")
	notExistLogger.Debug("you can see this debug log")
	notExistLogger.Info("you can see this info log")
	notExistLogger.Warn("you can see this warn log")
	notExistLogger.Error("you can see this error log")

	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[TRACE]-[notExist]-[test_virtual_logger.go:72] --- you can see this trace log\n"+
		"[DEBUG]-[notExist]-[test_virtual_logger.go:73] --- you can see this debug log\n"+
		"[INFO]-[notExist]-[test_virtual_logger.go:74] --- you can see this info log\n"+
		"[WARN]-[notExist]-[test_virtual_logger.go:75] --- you can see this warn log\n"+
		"[ERROR]-[notExist]-[test_virtual_logger.go:76] --- you can see this error log\n", content)

	time.Sleep(time.Second)
}
