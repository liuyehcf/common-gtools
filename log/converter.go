package log

import (
	"fmt"
	"strings"
)

const (
	leftAlign  = 0
	rightAlign = 1

	unlimitedWidth = -1
	blank          = ' '
)

type converter interface {
	// convert event to string
	convert(event *LoggingEvent) []byte

	// set next converter of converter chain
	setNext(next converter)

	// get next converter of converter chain
	getNext() converter
}

type abstractConverter struct {
	next      converter
	alignType int
	width     int
}

func (converter *abstractConverter) setNext(next converter) {
	converter.next = next
}

func (converter *abstractConverter) getNext() converter {
	return converter.next
}

func (converter *abstractConverter) truncAlign(content string) string {
	if converter.width == unlimitedWidth {
		return content
	}

	runes := []rune(content)
	if len(runes) >= converter.width {
		return content
	} else {
		extensionRunes := make([]rune, converter.width)
		if converter.alignType == leftAlign {
			copy(extensionRunes[:len(runes)], runes)
			copy(extensionRunes[len(runes):], fill(make([]rune, len(extensionRunes)-len(runes)), blank))
		} else {
			copy(extensionRunes[:converter.width-len(runes)], fill(make([]rune, converter.width-len(runes)), blank))
			copy(extensionRunes[converter.width-len(runes):], runes)
		}

		return string(extensionRunes)
	}
}

func fill(runes []rune, value rune) []rune {
	if runes == nil || len(runes) == 0 {
		return runes
	}

	runes[0] = value
	for i := 1; i < len(runes); i *= 2 {
		copy(runes[i:], runes[:i])
	}

	return runes
}

type headConverter struct {
	abstractConverter
}

func (converter *headConverter) convert(event *LoggingEvent) []byte {
	return []byte("")
}

// literal converter
type literalConverter struct {
	abstractConverter
	literal string
}

func (converter *literalConverter) convert(event *LoggingEvent) []byte {
	return []byte(converter.literal)
}

// logger converter
type loggerConverter struct {
	abstractConverter
}

func (converter *loggerConverter) convert(event *LoggingEvent) []byte {
	return []byte(converter.truncAlign(event.Name))
}

// date converter
type dateConverter struct {
	abstractConverter
	format string
}

func (converter *dateConverter) convert(event *LoggingEvent) []byte {
	return []byte(converter.truncAlign(event.Timestamp.Format(converter.format)))
}

// line converter
type lineConverter struct {
	abstractConverter
}

func (converter *lineConverter) convert(event *LoggingEvent) []byte {
	segments := strings.Split(event.File, pathSeparator)
	simpleFileName := segments[len(segments)-1]
	return []byte(converter.truncAlign(fmt.Sprintf("%s:%d", simpleFileName, event.Line)))
}

// message converter
type messageConverter struct {
	abstractConverter
}

func (converter *messageConverter) convert(event *LoggingEvent) []byte {
	return []byte(converter.truncAlign(event.GetFormattedMessage()))
}

// newline converter
type newlineConverter struct {
	abstractConverter
}

func (converter *newlineConverter) convert(event *LoggingEvent) []byte {
	return []byte("\n")
}

// level converter
type levelConverter struct {
	abstractConverter
}

func (converter *levelConverter) convert(event *LoggingEvent) []byte {
	switch event.Level {
	case TraceLevel:
		return []byte(converter.truncAlign("TRACE"))
	case DebugLevel:
		return []byte(converter.truncAlign("DEBUG"))
	case InfoLevel:
		return []byte(converter.truncAlign("INFO"))
	case WarnLevel:
		return []byte(converter.truncAlign("WARN"))
	case ErrorLevel:
		return []byte(converter.truncAlign("ERROR"))
	}

	panic(fmt.Sprintf("unsupported log level '%d'", event.Level))
}
