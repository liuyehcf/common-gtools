package test

import (
	"github.com/liuyehcf/common-gtools/utils"
	"regexp"
	"testing"
)

func TestRandomNumberString(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomNumberString(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomNumberString(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomNumberString(1)
		assertLength(s, 1)
		assertMatch("[0-9]", s)
	})
	execute(10000, func() {
		s := utils.RandomNumberString(6)
		assertLength(s, 6)
		assertMatch("[0-9]{6,6}", s)
	})
}

func TestRandomLetterStringLowerCase(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomLetterStringLowerCase(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringLowerCase(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringLowerCase(1)
		assertLength(s, 1)
		assertMatch("[a-z]", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringLowerCase(6)
		assertLength(s, 6)
		assertMatch("[a-z]{6,6}", s)
	})
}

func TestRandomLetterStringUpperCase(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomLetterStringUpperCase(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringUpperCase(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringUpperCase(1)
		assertLength(s, 1)
		assertMatch("[A-Z]", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterStringUpperCase(6)
		assertLength(s, 6)
		assertMatch("[A-Z]{6,6}", s)
	})
}

func TestRandomLetterNumberStringLowerCase(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomLetterNumberStringLowerCase(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringLowerCase(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringLowerCase(1)
		assertLength(s, 1)
		assertMatch("[a-z0-9]", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringLowerCase(6)
		assertLength(s, 6)
		assertMatch("[a-z0-9]{6,6}", s)
	})
}

func TestRandomLetterNumberStringUpperCase(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomLetterNumberStringUpperCase(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringUpperCase(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringUpperCase(1)
		assertLength(s, 1)
		assertMatch("[A-Z0-9]", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberStringUpperCase(6)
		assertLength(s, 6)
		assertMatch("[A-Z0-9]{6,6}", s)
	})
}

func TestRandomLetterNumberString(t *testing.T) {
	execute(10000, func() {
		s := utils.RandomLetterNumberString(-1)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberString(0)
		assertLength(s, 0)
		assertMatch("", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberString(1)
		assertLength(s, 1)
		assertMatch("[a-zA-Z0-9]", s)
	})
	execute(10000, func() {
		s := utils.RandomLetterNumberString(6)
		assertLength(s, 6)
		assertMatch("[a-zA-Z0-9]{6,6}", s)
	})
}

func execute(times int, f func()) {
	for i := 0; i < times; i += 1 {
		f()
	}
}

func assertLength(text string, length int) {
	utils.AssertTrue(len(text) == length, "test")
}

func assertMatch(pattern string, text string) {
	matched, err := regexp.MatchString(pattern, text)
	utils.AssertNil(err, "test")
	utils.AssertTrue(matched, "test")
}
