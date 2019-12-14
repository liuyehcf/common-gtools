package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/buffer"
	"github.com/liuyehcf/common-gtools/log"
	"time"
)

func main() {
	writer := log.NewStringWriter(buffer.NewRecycleByteBuffer(1024))
	writerAppender := log.NewWriterAppender(&log.AppenderConfig{
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
	assert.AssertTrue(content == "[INFO]-[additivityLogger]-[test_additivity.go:24] --- you can see this twice\n"+
		"[INFO]-[additivityLogger]-[test_additivity.go:24] --- you can see this twice\n", content)

	nonAdditivityLogger.Info("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[nonAdditivityLogger]-[test_additivity.go:30] --- you can see this once\n", content)
}
