package log

import (
	"github.com/liuyehcf/common-gtools/utils"
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
	// pre check is done here to avoid performing defer when the channel is closed, so as to improve performance
	if appender.isDestroyed {
		return
	}

	// although we made the flag bit judgment in the previous step,
	// this does not guarantee that the channel must be in normal state when sending data,
	// because it may be closed at this time.
	// That is to say, when checking the flag bit, the channel is normal. When sending data, the channel is closed
	//
	// recover from a crash caused by sending a message to a closed channel
	// it must be ensured that when the channel is closed, the flag 'isDestroyed' has been set to ensure the recovery of panic
	// so the isDestroyed flag must be set before the channel is closed
	defer appender.recoverIfChanClosed()
	if appender.filters == nil {
		// if channel is closed, then the upper statement will panic
		appender.queue <- appender.encoder.encode(event)
	} else {
		for _, filter := range appender.filters {
			if utils.IsNotNil(filter) {
				if !filter.Accept(event) {
					return
				}
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
