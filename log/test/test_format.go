package main

import (
	"github.com/liuyehcf/common-gtools/assert"
	"github.com/liuyehcf/common-gtools/log"
)

func main() {
	onePlaceHolder()
	twoPlaceHolder()
	multiply()
	escape()
	containSpace()
	chinese()
}

func onePlaceHolder() {
	assert.AssertTrue("abcdefg" == log.Format("{}", "abcdefg"), "test")
	assert.AssertTrue("{}" == log.Format("{}"), "test")
	assert.AssertTrue("abcdefg" == log.Format("{}", "abcdefg", "hijklmn"), "test")

	assert.AssertTrue("prefixabcdefg" == log.Format("prefix{}", "abcdefg"), "test")
	assert.AssertTrue("prefix{}" == log.Format("prefix{}"), "test")
	assert.AssertTrue("prefixabcdefg" == log.Format("prefix{}", "abcdefg", "hijklmn"), "test")

	assert.AssertTrue("abcdefgsuffix" == log.Format("{}suffix", "abcdefg"), "test")
	assert.AssertTrue("{}suffix" == log.Format("{}suffix"), "test")
	assert.AssertTrue("abcdefgsuffix" == log.Format("{}suffix", "abcdefg", "hijklmn"), "test")
}

func twoPlaceHolder() {
	assert.AssertTrue("a; b" == log.Format("{}; {}", "a", "b"), "test")
	assert.AssertTrue("{}; {}" == log.Format("{}; {}"), "test")
	assert.AssertTrue("a; {}" == log.Format("{}; {}", "a"), "test")
	assert.AssertTrue("a; b" == log.Format("{}; {}", "a", "b", "c"), "test")
	assert.AssertTrue("a; b" == log.Format("{}; {}", "a", "b", "c", "d"), "test")

	assert.AssertTrue("prefixa; b" == log.Format("prefix{}; {}", "a", "b"), "test")
	assert.AssertTrue("prefix{}; {}" == log.Format("prefix{}; {}"), "test")
	assert.AssertTrue("prefixa; {}" == log.Format("prefix{}; {}", "a"), "test")
	assert.AssertTrue("prefixa; b" == log.Format("prefix{}; {}", "a", "b", "c"), "test")
	assert.AssertTrue("prefixa; b" == log.Format("prefix{}; {}", "a", "b", "c", "d"), "test")

	assert.AssertTrue("a; bsuffix" == log.Format("{}; {}suffix", "a", "b"), "test")
	assert.AssertTrue("{}; {}suffix" == log.Format("{}; {}suffix"), "test")
	assert.AssertTrue("a; {}suffix" == log.Format("{}; {}suffix", "a"), "test")
	assert.AssertTrue("a; bsuffix" == log.Format("{}; {}suffix", "a", "b", "c"), "test")
	assert.AssertTrue("a; bsuffix" == log.Format("{}; {}suffix", "a", "b", "c", "d"), "test")
}

func multiply() {
	assert.AssertTrue("abcde" == log.Format("{}{}{}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("a bcde" == log.Format("{} {}{}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("ab cde" == log.Format("{}{} {}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("abc de" == log.Format("{}{}{} {}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("abcd e" == log.Format("{}{}{}{} {}", "a", "b", "c", "d", "e"), "test")
}

func escape() {
	assert.AssertTrue("\\{}" == log.Format("\\{}", "a"), "test")
	assert.AssertTrue("{\\}" == log.Format("{\\}", "a"), "test")
	assert.AssertTrue("\\{\\}" == log.Format("\\{\\}", "a"), "test")

	assert.AssertTrue("a; \\{}" == log.Format("{}; \\{}", "a"), "test")
	assert.AssertTrue("a; {\\}" == log.Format("{}; {\\}", "a"), "test")
	assert.AssertTrue("a; \\{\\}" == log.Format("{}; \\{\\}", "a"), "test")

	assert.AssertTrue("\\{}; a" == log.Format("\\{}; {}", "a"), "test")
	assert.AssertTrue("{\\}; a" == log.Format("{\\}; {}", "a"), "test")
	assert.AssertTrue("\\{\\}; a" == log.Format("\\{\\}; {}", "a"), "test")
}

func containSpace() {
	assert.AssertTrue("{ }" == log.Format("{ }", "a"), "test")
}

func chinese() {
	assert.AssertTrue("你好呀，小明" == log.Format("你好呀，{}", "小明"), "test")
}
