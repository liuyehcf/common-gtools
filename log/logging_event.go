package log

import (
	"time"
)

type LoggingEvent struct {
	Name             string
	Level            int
	Timestamp        time.Time
	File             string
	Line             int
	Message          string
	FormattedMessage string
	Values           []interface{}
	isInit           bool
}

func (event *LoggingEvent) GetFormattedMessage() string {
	if !event.isInit {
		event.FormattedMessage = format(event.Message, event.Values...)
		event.isInit = true
	}

	return event.FormattedMessage
}
