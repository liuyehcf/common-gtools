package log

import "fmt"

type Converter interface {
	// convert event to string
	Convert(event *LoggingEvent) []byte

	// set next converter of converter chain
	SetNext(next Converter)

	// get next converter of converter chain
	GetNext() Converter
}

type AbstractConverter struct {
	next Converter
}

func (converter *AbstractConverter) SetNext(next Converter) {
	converter.next = next
}

func (converter *AbstractConverter) GetNext() Converter {
	return converter.next
}

type HeadConverter struct {
	AbstractConverter
}

func (converter *HeadConverter) Convert(event *LoggingEvent) []byte {
	return []byte("")
}

// literal converter
type LiteralConverter struct {
	AbstractConverter
	literal string
}

func (converter *LiteralConverter) Convert(event *LoggingEvent) []byte {
	return []byte(converter.literal)
}

// level converter
type LevelConverter struct {
	AbstractConverter
}

func (converter *LevelConverter) Convert(event *LoggingEvent) []byte {
	switch event.Level {
	case TraceLevel:
		return []byte("TRACE")
	case DebugLevel:
		return []byte("DEBUG")
	case InfoLevel:
		return []byte("INFO")
	case WarnLevel:
		return []byte("WARN")
	case ErrorLevel:
		return []byte("ERROR")
	}
	panic(fmt.Sprintf("unsupported log level '%d'", event.Level))
}

// date converter
type DateConverter struct {
	AbstractConverter
	format string
}

func (converter *DateConverter) Convert(event *LoggingEvent) []byte {
	return []byte(event.Timestamp.Format(converter.format))
}

// message converter
type MessageConverter struct {
	AbstractConverter
}

func (converter *MessageConverter) Convert(event *LoggingEvent) []byte {
	return []byte(event.GetFormattedMessage())
}

// newline converter
type NewlineConverter struct {
	AbstractConverter
}

func (converter *NewlineConverter) Convert(event *LoggingEvent) []byte {
	return []byte("\n")
}
