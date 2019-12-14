package log

import (
	"io"
	"sync"
)

type AppenderConfig struct {
	// layout
	Layout string

	// filters of log
	Filters []Filter

	// only used for writerAppender
	Writer    io.WriteCloser
	NeedClose bool

	// only used for fileAppender
	FileRollingPolicy *RollingPolicy
}

type Appender interface {
	// do record this event
	DoAppend(event *LoggingEvent)

	// clean resources
	Destroy()
}

type abstractAppender struct {
	filters     []Filter
	encoder     encoder
	lock        *sync.Mutex
	queue       chan []byte
	isDestroyed bool
}

func (appender *abstractAppender) DoAppend(event *LoggingEvent) {
	if appender.isDestroyed {
		return
	}

	// recover from a crash caused by sending a message to a closed channel
	defer appender.recoverIfChanClosed()
	if appender.filters == nil {
		// if channel is closed, then the upper statement will panic
		appender.queue <- appender.encoder.encode(event)
	} else {
		for _, filter := range appender.filters {
			if !filter.Accept(event) {
				return
			}
		}
		// if channel is closed, then the upper statement will panic
		appender.queue <- appender.encoder.encode(event)
	}
}

func (appender *abstractAppender) recoverIfChanClosed() {
	if appender.isDestroyed {
		recover()
	}
}

func executeIgnorePanic(f func()) {
	defer func() {
		recover()
	}()
	f()
}
