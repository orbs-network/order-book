package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyParsing(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{in: "hello", expected: "hello"},
		{in: "a test", expected: "a_test"},
		{in: "aSnakeCaseString", expected: "a_snake_case_string"},
		{in: "multiple   spaces   are cleaned", expected: "multiple_spaces_are_cleaned"},
		{in: "ALL_UPPERCASE_IS_LOWERED", expected: "all_uppercase_is_lowered"},
		{in: "withMULTIPLEUppercaseIsBETTER", expected: "with_multiple_uppercase_is_better"},
		{in: "HTTPHandler", expected: "http_handler"},
		{in: "snake_case_is_preserved", expected: "snake_case_is_preserved"},
	}

	for _, tc := range cases {
		got := cleanKey(tc.in)
		assert.Equal(t, tc.expected, got)
	}
}
