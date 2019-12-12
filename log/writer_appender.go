package log

import (
	"github.com/liuyehcf/common-gtools/assert"
	"io"
	"sync"
)

type WriterAppender struct {
	AbstractAppender
	writer io.Writer
}

func NewWriterAppender(config *AppenderConfig) *WriterAppender {
	appender := &WriterAppender{
		AbstractAppender: AbstractAppender{
			encoder: NewPatternEncoder(config.Layout),
			filters: config.Filters,
			lock:    new(sync.Mutex),
			queue:   make(chan []byte, 1024),
		},
		writer: config.Writer,
	}

	go appender.onEventLoop()

	return appender
}

func (appender *WriterAppender) onEventLoop() {
	for {
		content := <-appender.queue
		appender.write(content)
	}
}

func (appender *WriterAppender) Destroy() {
	// do nothing
}

func (appender *WriterAppender) write(bytes []byte) {
	appender.lock.Lock()
	defer appender.lock.Unlock()
	_, err := appender.writer.Write(bytes)
	assert.AssertNil(err, "failed to write content")
}
