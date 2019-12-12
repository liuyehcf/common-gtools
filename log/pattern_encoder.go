package log

import (
	"bytes"
)

const (
	format = '%'
)

var (
	date    *Conversion
	message *Conversion
	newline *Conversion
	level   *Conversion
)

type Conversion struct {
	words []string
}

type PatternEncoder struct {
	layout string
	head   Converter
}

func NewPatternEncoder(layout string) *PatternEncoder {
	encoder := PatternEncoder{
		layout: layout,
	}

	encoder.initConverterChain()

	return &encoder
}

func (encoder *PatternEncoder) Encode(event *LoggingEvent) []byte {
	buffer := bytes.Buffer{}

	converter := encoder.head

	for ; converter != nil; {
		buffer.Write(converter.Convert(event))
		converter = converter.GetNext()
	}

	return buffer.Bytes()
}

func (encoder *PatternEncoder) initConverterChain() {
	runes := []rune(encoder.layout)

	index := 0
	encoder.head = &HeadConverter{}
	converter := encoder.head

	runesLen := len(runes)

	for ; index < runesLen; {
		c := runes[index]

		if c == format {
			if ok, offset := matchesConversion(runes, index+1, date); ok {
				index += offset + 1

				c = runes[index]
				if c != placeHolderStart {
					panic("unsupported date format '" + encoder.layout + "'")
				}

				buffer := bytes.Buffer{}
				index += 1
				if index < runesLen {
					c = runes[index]
				}

				for ; index < runesLen && c != placeHolderStop; {
					buffer.WriteRune(c)

					index += 1

					if index < runesLen {
						c = runes[index]
					}
				}

				if c != placeHolderStop {
					panic("unsupported date format '" + encoder.layout + "'")
				}

				nextConverter := &DateConverter{
					format: buffer.String(),
				}

				converter.SetNext(nextConverter)
				converter = nextConverter

				index += 1
			} else if ok, offset := matchesConversion(runes, index+1, level); ok {
				nextConverter := &LevelConverter{}
				converter.SetNext(nextConverter)
				converter = nextConverter

				index += offset + 1
			} else if ok, offset := matchesConversion(runes, index+1, message); ok {
				nextConverter := &MessageConverter{}
				converter.SetNext(nextConverter)
				converter = nextConverter

				index += offset + 1
			} else if ok, offset := matchesConversion(runes, index+1, newline); ok {
				nextConverter := &NewlineConverter{}
				converter.SetNext(nextConverter)
				converter = nextConverter

				index += offset + 1
			} else {
				panic("unsupported pattern '" + encoder.layout + "'")
			}
		} else {
			buffer := bytes.Buffer{}
			for ; index < runesLen && c != format; {
				buffer.WriteRune(c)

				index += 1

				if index < runesLen {
					c = runes[index]
				}
			}

			nextConverter := &LiteralConverter{
				literal: buffer.String(),
			}
			converter.SetNext(nextConverter)
			converter = nextConverter
		}
	}
}

func matchesConversion(runes []rune, start int, conversion *Conversion) (bool, int) {
	for _, word := range conversion.words {
		matches, offset := matchesWord(runes, start, word)

		if matches {
			return matches, offset
		}
	}

	return false, -1
}

func matchesWord(runes []rune, start int, word string) (bool, int) {
	expectedRunes := []rune(word)

	expectedLen := len(expectedRunes)

	if start+expectedLen > len(runes) {
		return false, -1
	}

	for i := 0; i < expectedLen; i += 1 {
		if runes[start+i] != expectedRunes[i] {
			return false, -1
		}
	}

	return true, expectedLen
}

func init() {
	date = &Conversion{
		words: []string{"d", "date"},
	}

	message = &Conversion{
		words: []string{"m", "msg", "message"},
	}

	newline = &Conversion{
		words: []string{"n"},
	}

	level = &Conversion{
		words: []string{"p", "le", "level"},
	}
}
