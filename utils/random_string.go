package utils

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	letterLowerCase = "abcdefghijklmnopqrstuvwxyz"
	letterUpperCase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number          = "0123456789"
)

func RandomNumberString(length int) string {
	return randomString([]rune(number), length)
}

func RandomLetterStringLowerCase(length int) string {
	return randomString([]rune(letterLowerCase), length)
}

func RandomLetterStringUpperCase(length int) string {
	return randomString([]rune(letterUpperCase), length)
}

func RandomLetterNumberStringLowerCase(length int) string {
	return randomString([]rune(letterLowerCase+number), length)
}

func RandomLetterNumberStringUpperCase(length int) string {
	return randomString([]rune(letterUpperCase+number), length)
}

func RandomLetterNumberString(length int) string {
	return randomString([]rune(letterLowerCase+letterUpperCase+number), length)
}

func randomString(runes []rune, length int) string {
	buffer := bytes.Buffer{}

	size := len(runes)

	for i := 0; i < length; i += 1 {
		buffer.WriteRune(runes[rand.Intn(size)])
	}

	return buffer.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
