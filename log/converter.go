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

type Converter interface {
	// convert event to string
	Convert(event *LoggingEvent) []byte

	// set next converter of converter chain
	SetNext(next Converter)

	// get next converter of converter chain
	GetNext() Converter
}

type AbstractConverter struct {
	next      Converter
	alignType int
	width     int
}

func (converter *AbstractConverter) SetNext(next Converter) {
	converter.next = next
}

func (converter *AbstractConverter) GetNext() Converter {
	return converter.next
}

func (converter *AbstractConverter) truncAlign(content string) string {
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

			for i := len(runes); i < len(extensionRunes); i += 1 {
				extensionRunes[i] = blank
			}
		} else {
			for i := 0; i < converter.width-len(runes); i += 1 {
				extensionRunes[i] = blank
			}

			copy(extensionRunes[converter.width-len(runes):], runes)
		}

		return string(extensionRunes)
	}
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

// date converter
type DateConverter struct {
	AbstractConverter
	format string
}

func (converter *DateConverter) Convert(event *LoggingEvent) []byte {
	return []byte(converter.truncAlign(event.Timestamp.Format(converter.format)))
}

// line converter
type LineConverter struct {
	AbstractConverter
}

func (converter *LineConverter) Convert(event *LoggingEvent) []byte {
	segments := strings.Split(event.File, pathSeparator)
	simpleFileName := segments[len(segments)-1]
	return []byte(converter.truncAlign(fmt.Sprintf("%s:%d", simpleFileName, event.Line)))
}

// message converter
type MessageConverter struct {
	AbstractConverter
}

func (converter *MessageConverter) Convert(event *LoggingEvent) []byte {
	return []byte(converter.truncAlign(event.GetFormattedMessage()))
}

// newline converter
type NewlineConverter struct {
	AbstractConverter
}

func (converter *NewlineConverter) Convert(event *LoggingEvent) []byte {
	return []byte("\n")
}

// level converter
type LevelConverter struct {
	AbstractConverter
}

func (converter *LevelConverter) Convert(event *LoggingEvent) []byte {
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