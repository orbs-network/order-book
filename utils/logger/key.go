package logger

import (
	"strings"
	"unicode"
)

// cleanKey will force all keys to snakecase.
func cleanKey(key string) string {
	out := ""
	for i := range key {
		switch {
		case isSpace(key[i]) && i > 0 && isSpace(key[i-1]):
			// ignore multiple spaces - collapse to 1
			// the order is important here - it must occur before
			// the test for a whitespace char.
		case isSpace(key[i]):
			out += "_"
		case isUpper(key[i]) && i > 0 && isLower(key[i-1]):
			out += "_" + lower(key[i])
		case isUpper(key[i]) && i > 0 && isUpper(key[i-1]) && i < len(key)-1 && isLower(key[i+1]):
			out += "_" + lower(key[i])
		default:
			out += lower(key[i])
		}
	}

	return out
}

func isUpper(c byte) bool {
	return unicode.IsUpper(rune(c))
}

func isLower(c byte) bool {
	return unicode.IsLower(rune(c))
}

func isSpace(c byte) bool {
	return string(c) == " "
}

func lower(c byte) string {
	return strings.ToLower(string(c))
}
