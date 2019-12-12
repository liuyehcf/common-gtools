package log

import (
	"bytes"
	"github.com/liuyehcf/common-gtools/assert"
	"strconv"
)

const (
	percent       = '%'
	strikethrough = '-'
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

		if c == percent {
			index += 1

			// parse align
			alignType, offset := getAlignType(runes, index)
			index += offset

			// parse width
			width, offset := getWidth(runes, index)
			index += offset

			if ok, offset := matchesConversion(runes, index, date); ok {
				index += offset

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
					AbstractConverter: AbstractConverter{
						alignType: alignType,
						width:     width,
					},
					format: buffer.String(),
				}

				converter.SetNext(nextConverter)
				converter = nextConverter

				index += 1
			} else if ok, offset := matchesConversion(runes, index, level); ok {
				index += offset
				nextConverter := &LevelConverter{
					AbstractConverter: AbstractConverter{
						alignType: alignType,
						width:     width,
					},
				}
				converter.SetNext(nextConverter)
				converter = nextConverter
			} else if ok, offset := matchesConversion(runes, index, message); ok {
				index += offset
				nextConverter := &MessageConverter{
					AbstractConverter: AbstractConverter{
						alignType: alignType,
						width:     width,
					},
				}
				converter.SetNext(nextConverter)
				converter = nextConverter
			} else if ok, offset := matchesConversion(runes, index, newline); ok {
				index += offset
				nextConverter := &NewlineConverter{
					AbstractConverter: AbstractConverter{
						alignType: alignType,
						width:     width,
					},
				}
				converter.SetNext(nextConverter)
				converter = nextConverter
			} else {
				panic("unsupported pattern '" + encoder.layout + "'")
			}
		} else {
			buffer := bytes.Buffer{}
			for ; index < runesLen && c != percent; {
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

func getAlignType(runes []rune, start int) (int, int) {
	if start >= len(runes) {
		return rightAlign, 0
	}

	if runes[start] == strikethrough {
		return leftAlign, 1
	} else {
		return rightAlign, 0
	}
}

func getWidth(runes []rune, start int) (int, int) {
	if start >= len(runes) {
		return unlimitedWidth, 0
	}

	index := start

	for ; index < len(runes); {
		v := runes[index]
		if v < '0' || v > '9' {
			break
		}
		index += 1
	}

	if index == start {
		return unlimitedWidth, 0
	}

	width, err := strconv.Atoi(string(runes[start:index]))
	assert.AssertNil(err, "failed to parse width")

	return width, index - start
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
