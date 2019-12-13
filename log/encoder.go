package log

type encoder interface {
	// encoding logging event to bytes
	encode(event *LoggingEvent) []byte
}
