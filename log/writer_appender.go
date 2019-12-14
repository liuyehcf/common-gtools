package log

import (
	"github.com/liuyehcf/common-gtools/assert"
	"io"
	"sync"
)

type writerAppender struct {
	abstractAppender
	writer    io.WriteCloser
	needClose bool
}

func NewWriterAppender(config *AppenderConfig) *writerAppender {
	assert.AssertNotNil(config.Writer, "write is required for writer appender")

	appender := &writerAppender{
		abstractAppender: abstractAppender{
			encoder: newPatternEncoder(config.Layout),
			filters: config.Filters,
			lock:    new(sync.Mutex),
			queue:   make(chan []byte, 1024),
		},
		writer:    config.Writer,
		needClose: config.NeedClose,
	}

	go appender.onEventLoop()

	return appender
}

func (appender *writerAppender) Destroy() {
	lock.Lock()
	defer lock.Unlock()
	appender.isDestroyed = true
	executeIgnorePanic(func() {
		close(appender.queue)
	})
	executeIgnorePanic(func() {
		if appender.needClose {
			_ = appender.writer.Close()
		}
	})
}

func (appender *writerAppender) onEventLoop() {
	defer func() {
		recover()
	}()

	var content []byte
	var ok bool
	for !appender.isDestroyed {
		if content, ok = <-appender.queue; !ok {
			// channel is closed
			break

		}
		appender.write(content)
	}
}

func (appender *writerAppender) write(bytes []byte) {
	appender.lock.Lock()
	defer appender.lock.Unlock()
	_, err := appender.writer.Write(bytes)
	assert.AssertNil(err, "failed to write content")
}
