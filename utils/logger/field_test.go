package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerFields(t *testing.T) {
	cases := []struct {
		name     string
		input    Field
		expected zap.Field
	}{
		{name: "string field", input: String("key", "value"), expected: zap.Field{Key: "key", String: "value", Type: zapcore.StringType}},
		{name: "error field", input: Error(errors.New("an error")), expected: zap.Field{Key: "error", Interface: errors.New("an error"), Type: zapcore.ErrorType}},
		{name: "bool field", input: Bool("key", true), expected: zap.Field{Key: "key", Integer: 1, Type: zapcore.BoolType}},
		{name: "int field", input: Int("key", 100), expected: zap.Field{Key: "key", Integer: 100, Type: zapcore.Int64Type}},
		{name: "int8 field", input: Int8("key", 7), expected: zap.Field{Key: "key", Integer: 7, Type: zapcore.Int8Type}},
		{name: "int16 field", input: Int16("key", 2), expected: zap.Field{Key: "key", Integer: 2, Type: zapcore.Int16Type}},
		{name: "int32 field", input: Int32("key", 1001), expected: zap.Field{Key: "key", Integer: 1001, Type: zapcore.Int32Type}},
		{name: "int64 field", input: Int64("key", 72), expected: zap.Field{Key: "key", Integer: 72, Type: zapcore.Int64Type}},
		{name: "int field", input: Uint("key", 100), expected: zap.Field{Key: "key", Integer: 100, Type: zapcore.Uint64Type}},
		{name: "int8 field", input: Uint8("key", 7), expected: zap.Field{Key: "key", Integer: 7, Type: zapcore.Uint8Type}},
		{name: "int16 field", input: Uint16("key", 2), expected: zap.Field{Key: "key", Integer: 2, Type: zapcore.Uint16Type}},
		{name: "int32 field", input: Uint32("key", 1001), expected: zap.Field{Key: "key", Integer: 1001, Type: zapcore.Uint32Type}},
		{name: "int64 field", input: Uint64("key", 72), expected: zap.Field{Key: "key", Integer: 72, Type: zapcore.Uint64Type}},
		{name: "float32 field", input: Float32("key", 12.2), expected: zap.Field{Key: "key", Integer: 1094923059, Type: zapcore.Float32Type}},
		{name: "float64 field", input: Float64("key", 90.6), expected: zap.Field{Key: "key", Integer: 4636075825159366246, Type: zapcore.Float64Type}},
		{name: "bytes field", input: Bytes("key", []byte("some bytes")), expected: zap.Field{Key: "key", Interface: []byte("some bytes"), Type: zapcore.BinaryType}},
		{name: "byte string field", input: ByteString("key", []byte("some bytes")), expected: zap.Field{Key: "key", Interface: []byte("some bytes"), Type: zapcore.ByteStringType}},
	}
	for _, tc := range cases {
		t.Run("string method", func(t *testing.T) {
			l, o := newTestLogger()

			l.Info("message", tc.input)

			require.Equal(t, 1, o.Len())
			logs := o.All()
			assert.Equal(t, "message", logs[0].Message)
			assert.ElementsMatch(t, []zap.Field{tc.expected}, logs[0].Context)
		})
	}
}

func assertFieldIn(t *testing.T, fields map[string]Field, expected Field) {
	f, ok := fields[expected.Key()]
	require.True(t, ok)

	assert.Equal(t, expected.ToZap()[0], f.ToZap()[0])
}

func assertFieldNotIn(t *testing.T, fields map[string]Field, expected Field) {
	_, ok := fields[expected.Key()]
	require.False(t, ok)
}
