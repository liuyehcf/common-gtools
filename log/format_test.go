package log

import (
	"github.com/liuyehcf/common-gtools/utils"
	"testing"
)

func TestOnePlaceHolder(t *testing.T) {
	utils.AssertTrue("abcdefg" == format("{}", "abcdefg"), "test")
	utils.AssertTrue("{}" == format("{}"), "test")
	utils.AssertTrue("abcdefg" == format("{}", "abcdefg", "hijklmn"), "test")

	utils.AssertTrue("prefixabcdefg" == format("prefix{}", "abcdefg"), "test")
	utils.AssertTrue("prefix{}" == format("prefix{}"), "test")
	utils.AssertTrue("prefixabcdefg" == format("prefix{}", "abcdefg", "hijklmn"), "test")

	utils.AssertTrue("abcdefgsuffix" == format("{}suffix", "abcdefg"), "test")
	utils.AssertTrue("{}suffix" == format("{}suffix"), "test")
	utils.AssertTrue("abcdefgsuffix" == format("{}suffix", "abcdefg", "hijklmn"), "test")
}

func TestTwoPlaceHolder(t *testing.T) {
	utils.AssertTrue("a; b" == format("{}; {}", "a", "b"), "test")
	utils.AssertTrue("{}; {}" == format("{}; {}"), "test")
	utils.AssertTrue("a; {}" == format("{}; {}", "a"), "test")
	utils.AssertTrue("a; b" == format("{}; {}", "a", "b", "c"), "test")
	utils.AssertTrue("a; b" == format("{}; {}", "a", "b", "c", "d"), "test")

	utils.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b"), "test")
	utils.AssertTrue("prefix{}; {}" == format("prefix{}; {}"), "test")
	utils.AssertTrue("prefixa; {}" == format("prefix{}; {}", "a"), "test")
	utils.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b", "c"), "test")
	utils.AssertTrue("prefixa; b" == format("prefix{}; {}", "a", "b", "c", "d"), "test")

	utils.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b"), "test")
	utils.AssertTrue("{}; {}suffix" == format("{}; {}suffix"), "test")
	utils.AssertTrue("a; {}suffix" == format("{}; {}suffix", "a"), "test")
	utils.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b", "c"), "test")
	utils.AssertTrue("a; bsuffix" == format("{}; {}suffix", "a", "b", "c", "d"), "test")
}

func TestMultiply(t *testing.T) {
	utils.AssertTrue("abcde" == format("{}{}{}{}{}", "a", "b", "c", "d", "e"), "test")
	utils.AssertTrue("a bcde" == format("{} {}{}{}{}", "a", "b", "c", "d", "e"), "test")
	utils.AssertTrue("ab cde" == format("{}{} {}{}{}", "a", "b", "c", "d", "e"), "test")
	utils.AssertTrue("abc de" == format("{}{}{} {}{}", "a", "b", "c", "d", "e"), "test")
	utils.AssertTrue("abcd e" == format("{}{}{}{} {}", "a", "b", "c", "d", "e"), "test")
}

func TestEscape(t *testing.T) {
	utils.AssertTrue("\\{}" == format("\\{}", "a"), "test")
	utils.AssertTrue("{\\}" == format("{\\}", "a"), "test")
	utils.AssertTrue("\\{\\}" == format("\\{\\}", "a"), "test")

	utils.AssertTrue("a; \\{}" == format("{}; \\{}", "a"), "test")
	utils.AssertTrue("a; {\\}" == format("{}; {\\}", "a"), "test")
	utils.AssertTrue("a; \\{\\}" == format("{}; \\{\\}", "a"), "test")

	utils.AssertTrue("\\{}; a" == format("\\{}; {}", "a"), "test")
	utils.AssertTrue("{\\}; a" == format("{\\}; {}", "a"), "test")
	utils.AssertTrue("\\{\\}; a" == format("\\{\\}; {}", "a"), "test")
}

func TestContainSpace(t *testing.T) {
	utils.AssertTrue("{ }" == format("{ }", "a"), "test")
}

func TestChinese(t *testing.T) {
	utils.AssertTrue("你好呀，小明" == format("你好呀，{}", "小明"), "test")
}
