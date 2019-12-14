package log

import (
	"github.com/liuyehcf/common-gtools/assert"
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
	assert.AssertTrue("abcdefg" == format("{}", "abcdefg"), "test")
	assert.AssertTrue("{}" == format("{}"), "test")
	assert.AssertTrue("abcdefg" == format("{}", "abcdefg", "hijklmn"), "test")

	assert.AssertTrue("prefixabcdefg" == format("prefix{}", "abcdefg"), "test")
	assert.AssertTrue("prefix{}" == format("prefix{}"), "test")
	assert.AssertTrue("prefixabcdefg" == format("prefix{}", "abcdefg", "hijklmn"), "test")

	assert.AssertTrue("abcdefgsuffix" == format("{}suffix", "abcdefg"), "test")
	assert.AssertTrue("{}suffix" == format("{}suffix"), "test")
	assert.AssertTrue("abcdefgsuffix" == format("{}suffix", "abcdefg", "hijklmn"), "test")
}

func twoPlaceHolder() {
	assert.AssertTrue("a; b" == format("{}; {}", "a", "b"), "test")
	assert.AssertTrue("{}; {}" == format("{}; {}"), "test")
	assert.AssertTrue("a; {}" == format("{}; {}", "a"), "test")
	assert.AssertTrue("a; b" == format("{}; {}", "a", "b", "c"), "test")
	assert.AssertTrue("a; b" == format("{}; {}", "a", "b", "c", "d"), "test")

	assert.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b"), "test")
	assert.AssertTrue("prefix{}; {}" == format("prefix{}; {}"), "test")
	assert.AssertTrue("prefixa; {}" == format("prefix{}; {}", "a"), "test")
	assert.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b", "c"), "test")
	assert.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b", "c", "d"), "test")

	assert.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b"), "test")
	assert.AssertTrue("{}; {}suffix" == format("{}; {}suffix"), "test")
	assert.AssertTrue("a; {}suffix" == format("{}; {}suffix", "a"), "test")
	assert.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b", "c"), "test")
	assert.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b", "c", "d"), "test")
}

func multiply() {
	assert.AssertTrue("abcde" == format("{}{}{}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("a bcde" == format("{} {}{}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("ab cde" == format("{}{} {}{}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("abc de" == format("{}{}{} {}{}", "a", "b", "c", "d", "e"), "test")
	assert.AssertTrue("abcd e" == format("{}{}{}{} {}", "a", "b", "c", "d", "e"), "test")
}

func escape() {
	assert.AssertTrue("\\{}" == format("\\{}", "a"), "test")
	assert.AssertTrue("{\\}" == format("{\\}", "a"), "test")
	assert.AssertTrue("\\{\\}" == format("\\{\\}", "a"), "test")

	assert.AssertTrue("a; \\{}" == format("{}; \\{}", "a"), "test")
	assert.AssertTrue("a; {\\}" == format("{}; {\\}", "a"), "test")
	assert.AssertTrue("a; \\{\\}" == format("{}; \\{\\}", "a"), "test")

	assert.AssertTrue("\\{}; a" == format("\\{}; {}", "a"), "test")
	assert.AssertTrue("{\\}; a" == format("{\\}; {}", "a"), "test")
	assert.AssertTrue("\\{\\}; a" == format("\\{\\}; {}", "a"), "test")
}

func containSpace() {
	assert.AssertTrue("{ }" == format("{ }", "a"), "test")
}

func chinese() {
	assert.AssertTrue("你好呀，小明" == format("你好呀，{}", "小明"), "test")
}
