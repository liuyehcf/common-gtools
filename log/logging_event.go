package log

import (
	"time"
)

type LoggingEvent struct {
	Level            int
	Timestamp        time.Time
	Message          string
	FormattedMessage string
	Values           []interface{}
	isInit           bool
}

func (event *LoggingEvent) GetFormattedMessage() string {
	if !event.isInit {
		event.FormattedMessage = Format(event.Message, event.Values...)
		event.isInit = true
	}

	return event.FormattedMessage
}
