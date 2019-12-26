package log

import (
	"errors"
	"io"
	"sync"
)

type writerAppender struct {
	abstractAppender
	writer    io.WriteCloser
	needClose bool
}

func NewWriterAppender(config *AppenderConfig) (*writerAppender, error) {
	if config.Writer == nil {
		return nil, errors.New("write is required for writer appender")
	}

	encoder, err := newPatternEncoder(config.Layout)
	if err != nil {
		return nil, err
	}
	appender := &writerAppender{
		abstractAppender: abstractAppender{
			encoder: encoder,
			filters: config.Filters,
			lock:    new(sync.Mutex),
			queue:   make(chan []byte, 1024),
		},
		writer:    config.Writer,
		needClose: config.NeedClose,
	}

	go appender.onEventLoop()

	return appender, nil
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
	_, _ = appender.writer.Write(bytes)
}
