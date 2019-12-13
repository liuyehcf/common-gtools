package log

import (
	"github.com/liuyehcf/common-gtools/assert"
	"io"
	"sync"
)

type writerAppender struct {
	abstractAppender
	writer io.Writer
}

func NewWriterAppender(config *AppenderConfig) *writerAppender {
	appender := &writerAppender{
		abstractAppender: abstractAppender{
			encoder: newPatternEncoder(config.Layout),
			filters: config.Filters,
			lock:    new(sync.Mutex),
			queue:   make(chan []byte, 1024),
		},
		writer: config.Writer,
	}

	go appender.onEventLoop()

	return appender
}

func (appender *writerAppender) onEventLoop() {
	for {
		content := <-appender.queue
		appender.write(content)
	}
}

func (appender *writerAppender) Destroy() {
	close(appender.queue)
}

func (appender *writerAppender) write(bytes []byte) {
	appender.lock.Lock()
	defer appender.lock.Unlock()
	_, err := appender.writer.Write(bytes)
	assert.AssertNil(err, "failed to write content")
}
