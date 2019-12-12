package log

type Encoder interface {
	// encoding logging event to bytes
	Encode(event *LoggingEvent) []byte
}
