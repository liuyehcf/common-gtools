package log

import (
	"io"
	"sync"
)

type AppenderConfig struct {
	Layout  string
	Filters []Filter

	// only used for WriterAppender
	Writer io.Writer

	// only used for FileAppender
	FileRollingPolicy *RollingPolicy
}

type Appender interface {
	// do record this event
	DoAppend(event *LoggingEvent)

	// clean resources
	Destroy()
}

type AbstractAppender struct {
	filters []Filter
	encoder Encoder
	lock    *sync.Mutex
	queue   chan []byte
}

func (appender *AbstractAppender) DoAppend(event *LoggingEvent) {
	if appender.filters == nil {
		appender.queue <- appender.encoder.Encode(event)
	} else {
		for _, filter := range appender.filters {
			if !filter.Accept(event) {
				return
			}
		}

		appender.queue <- appender.encoder.Encode(event)
	}
}
