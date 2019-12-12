package log

type Filter interface {
	// whether accept specified log event
	Accept(event *LoggingEvent) bool
}

type LevelFilter struct {
	LogLevelThreshold int
}

func (filter *LevelFilter) Accept(event *LoggingEvent) bool {
	return event.Level >= filter.LogLevelThreshold
}
