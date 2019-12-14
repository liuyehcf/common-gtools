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

	logger := log.NewLogger("non-root", log.InfoLevel, false, []log.Appender{writerAppender})

	var content string

	logger.Info("you can see this twice")
	logger.Error("you can see this twice")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[INFO]-[non-root]-[test_change_non_root_logger.go:22] --- you can see this twice\n"+
		"[ERROR]-[non-root]-[test_change_non_root_logger.go:23] --- you can see this twice\n", content)

	log.NewLogger("non-root", log.ErrorLevel, false, []log.Appender{writerAppender})

	logger.Info("you cannot see this")
	logger.Error("you can see this once")
	time.Sleep(time.Millisecond * 10)
	content = writer.ReadString()
	assert.AssertTrue(content == "[ERROR]-[non-root]-[test_change_non_root_logger.go:32] --- you can see this once\n", content)
}
