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
	filters []Filter
	encoder encoder
	lock    *sync.Mutex
	queue   chan []byte
}

func (appender *abstractAppender) DoAppend(event *LoggingEvent) {
	if appender.filters == nil {
		appender.queue <- appender.encoder.encode(event)
	} else {
		for _, filter := range appender.filters {
			if !filter.Accept(event) {
				return
			}
		}

		appender.queue <- appender.encoder.encode(event)
	}
}
