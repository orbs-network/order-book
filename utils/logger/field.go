package logger

import (
	"strconv"

	"go.uber.org/zap"
)

// Field represents the different field types that can be passed in.
type Field interface {
	Key() string
	ToZap() []zap.Field
}

func mapToZap(defaultFields map[string]Field, fields []Field) []zap.Field {
	combined := append(toSlice(defaultFields), fields...)
	zapFields := make([]zap.Field, 0, len(combined))
	for i := range combined {
		zapFields = append(zapFields, combined[i].ToZap()...)
	}

	return zapFields
}

func toSlice(in map[string]Field) []Field {
	out := make([]Field, 0)
	for _, f := range in {
		out = append(out, f)
	}

	return out
}

// Error creates a field for logging error types.
func Error(err error) FieldWrapper[error] {
	return NewWrapper("error", err, zap.NamedError)
}

// String creates a field for logging string types.
func String(key string, val string) FieldWrapper[string] {
	return NewWrapper(key, val, zap.String)
}

// Strings creates a field for logging string types.
func Strings(key string, val []string) FieldWrapper[[]string] {
	return NewWrapper(key, val, zap.Strings)
}

// Bool creates a field for logging bool types.
func Bool(key string, val bool) FieldWrapper[bool] {
	return NewWrapper(key, val, zap.Bool)
}

// Uint creates a field for logging uint types.
func Uint(key string, val uint) FieldWrapper[uint] {
	return NewWrapper(key, val, zap.Uint)
}

// Uint8 creates a field for logging uint8 types.
func Uint8(key string, val uint8) FieldWrapper[uint8] {
	return NewWrapper(key, val, zap.Uint8)
}

// Uint16 creates a field for logging uint16 types.
func Uint16(key string, val uint16) FieldWrapper[uint16] {
	return NewWrapper(key, val, zap.Uint16)
}

// Uint32 creates a field for logging uint32 types.
func Uint32(key string, val uint32) FieldWrapper[uint32] {
	return NewWrapper(key, val, zap.Uint32)
}

// Uint64 creates a field for logging uint64 types.
func Uint64(key string, val uint64) FieldWrapper[uint64] {
	return NewWrapper(key, val, zap.Uint64)
}

// Int creates a field for logging int types.
func Int(key string, val int) FieldWrapper[int] {
	return NewWrapper(key, val, zap.Int)
}

// Int8 creates a field for logging int8 types.
func Int8(key string, val int8) FieldWrapper[int8] {
	return NewWrapper(key, val, zap.Int8)
}

// Int16 creates a field for logging int16 types.
func Int16(key string, val int16) FieldWrapper[int16] {
	return NewWrapper(key, val, zap.Int16)
}

// Int32 creates a field for logging int32 types.
func Int32(key string, val int32) FieldWrapper[int32] {
	return NewWrapper(key, val, zap.Int32)
}

// Int64 creates a field for logging int64 types.
func Int64(key string, val int64) FieldWrapper[int64] {
	return NewWrapper(key, val, zap.Int64)
}

// Float32 creates a field for logging float32 types.
func Float32(key string, val float32) FieldWrapper[float32] {
	return NewWrapper(key, val, zap.Float32)
}

// Float64 creates a field for logging float64 types.
func Float64(key string, val float64) FieldWrapper[float64] {
	return NewWrapper(key, val, zap.Float64)
}

// Bytes creates a field for logging an array of byte types.
func Bytes(key string, val []byte) FieldWrapper[[]byte] {
	return NewWrapper(key, val, zap.Binary)
}

// ByteString creates a field for logging valid utf-8 byte arrays to a string representation.
func ByteString(key string, val []byte) FieldWrapper[[]byte] {
	return NewWrapper(key, val, zap.ByteString)
}

// DDSpanID creates a field for logging datadog span IDs in a readable format.
func DDSpanID(val uint64) FieldWrapper[string] {
	return NewWrapper("dd.span_id", strconv.FormatUint(val, 10), zap.String)
}

// DDTraceID creates a field for logging datadog trace IDs in a readable format.
func DDTraceID(val uint64) FieldWrapper[string] {
	return NewWrapper("dd.trace_id", strconv.FormatUint(val, 10), zap.String)
}

// NewWrapper is a helper method for creating FieldWrapper.
func NewWrapper[T any](key string, val T, toZap func(key string, val T) zap.Field) FieldWrapper[T] {
	return FieldWrapper[T]{
		key:   cleanKey(key),
		val:   val,
		toZap: toZap,
	}
}

// FieldWrapper is a generic struct abstracting the field for creating concrete implementations
// for Fields.
//
// This implementation currently only supports one to one mapping between fields.
type FieldWrapper[T any] struct {
	key   string
	val   T
	toZap func(key string, val T) zap.Field
}

// ToZap implements the Field interface for writing log fields to zap types.
func (f FieldWrapper[T]) ToZap() []zap.Field {
	return []zap.Field{f.toZap(f.key, f.val)}
}

// Key returns the fields key.
func (f FieldWrapper[T]) Key() string {
	return f.key
}
