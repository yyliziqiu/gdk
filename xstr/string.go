package xstr

import (
	"math/rand"
	"strings"
	"unicode/utf8"
)

// Truncate 截断字符串
func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// TruncateUtf8 截断 UTF8 编码的字符串
func TruncateUtf8(str string, n int) string {
	if len(str) <= n {
		return str
	}

	if utf8.RuneCountInString(str) <= n {
		return str
	}

	return string(([]rune(str))[:n])
}

// TrimSplit 分割字符串并去除分割结果两侧的空格
func TrimSplit(str string, sep string) []string {
	ret := strings.Split(str, sep)
	for i := 0; i < len(ret); i++ {
		ret[i] = strings.TrimSpace(ret[i])
	}
	return ret
}

const _randomDigit = "0123456789"

// RandomDigit 随机只包含数字的字符串
func RandomDigit(length int) string {
	return Random(_randomDigit, length)
}

const _randomAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomAlphabet 随机只包含字母的字符串
func RandomAlphabet(length int) string {
	return Random(_randomAlphabet, length)
}

const _randomCharsets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString 随机包含数字和字母的字符串
func RandomString(length int) string {
	return Random(_randomCharsets, length)
}

// Random 随机字符串
func Random(charset string, length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(62)])
	}

	return sb.String()
}
