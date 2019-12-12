package log

import (
	"bytes"
	"fmt"
)

const (
	placeHolderStart = '{'
	placeHolderStop  = '}'
	escapeChar       = '\\'
)

func Format(format string, values ...interface{}) string {
	if values == nil || len(values) == 0 {
		return format
	}
	valueLen := len(values)

	var buffer bytes.Buffer

	runeFormat := []rune(format)
	placeHolderIndex := 0

	isCurEscapeChar := false
	isPreEscapeChar := false

	isCurPlaceHolderStart := false
	isPrePlaceHolderStart := false

	isNoMoreValue := false

	for _, c := range runeFormat {
		isCurEscapeChar = false
		isCurPlaceHolderStart = false

		needAppend := true

		if c == placeHolderStart {
			if !isPreEscapeChar {
				isCurPlaceHolderStart = true
			}
		} else if c == placeHolderStop {
			// if there is no more value to replace the remaining placeholders, we just keep the original string
			if !isNoMoreValue && isPrePlaceHolderStart {
				needAppend = false

				// remove the pre '{'
				buffer.Truncate(buffer.Len() - 1)
				value := values[placeHolderIndex]
				buffer.WriteString(stringify(value))

				placeHolderIndex += 1
				if placeHolderIndex >= valueLen {
					isNoMoreValue = true
				}
			}
		} else if c == escapeChar {
			if !isPreEscapeChar {
				isCurEscapeChar = true
			}
		}

		if needAppend {
			buffer.WriteRune(c)
		}

		isPreEscapeChar = isCurEscapeChar
		isPrePlaceHolderStart = isCurPlaceHolderStart
	}

	return buffer.String()
}

func stringify(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
