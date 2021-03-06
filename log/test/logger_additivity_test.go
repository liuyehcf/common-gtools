package main

import (
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"github.com/liuyehcf/common-gtools/utils"
	"testing"
	"time"
)

func TestAdditivity(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	additivityLogger := log.NewLogger("additivityLogger", log.InfoLevel, true, []log.Appender{writerAppender})
	nonAdditivityLogger := log.NewLogger("nonAdditivityLogger", log.InfoLevel, false, []log.Appender{writerAppender})
	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	additivityLogger.Info("you can see this twice")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "[INFO]-[additivityLogger]-[logger_additivity_test.go:25] --- you can see this twice\n"+
		"[INFO]-[additivityLogger]-[logger_additivity_test.go:25] --- you can see this twice\n", content)

	nonAdditivityLogger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "[INFO]-[nonAdditivityLogger]-[logger_additivity_test.go:31] --- you can see this once\n", content)
}

func TestChangeRootLogger(t *testing.T) {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender1, _ := log.NewWriterAppender(&log.AppenderConfig{
		Layout:  "[%p]-[%c]-[%L] --- %m%n",
		Filters: nil,
		Writer:  writer,
	})

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{writerAppender1})
	logger := log.GetLogger("notExist")

	var content string

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "[INFO]-[notExist]-[logger_additivity_test.go:50] --- you can see this once\n", content)

	log.NewLogger(log.Root, log.InfoLevel, false, []log.Appender{nil})
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "[WARN]-[ROOT]-[logger.go:191] --- logger 'ROOT' contains nil appender\n"+
		"[WARN]-[ROOT]-[logger.go:224] --- logger 'ROOT' is replaced\n", content)

	logger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	utils.AssertTrue(content == "", content)
}
