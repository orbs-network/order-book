// Package logctx is a small package designed for propagating a
// production-ready logger through contexts.
package logctx

import (
	"context"

	"github.com/orbs-network/order-book/utils/logger"
)

// defaultLogger is a production-ready logger. From will use defaultLogger
// if there is no logger already present in the context.
var defaultLogger = logger.New(logger.WithCallerSkip(3))

// key can only ever be one value and will not allocate when doing lookups.
// See https://github.com/golang/go/issues/17826.
type key struct{}

// From will return the logger associated with the context if present,
// otherwise it will return Default.
func From(ctx context.Context) logger.Logger {
	if l, ok := ctx.Value(key{}).(logger.Logger); ok {
		return l
	}
	return defaultLogger
}

// With will add the logger to the context.
func With(ctx context.Context, l logger.Logger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

// WithFields will return a new context with the fields added to the associated
// logger.
func WithFields(ctx context.Context, fields ...logger.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	return With(ctx, logger.WithFields(From(ctx), fields...))
}

// Info will log at the info level using the logger associated with the
// context.
func Info(ctx context.Context, msg string, fields ...logger.Field) {
	From(ctx).Info(msg, fields...)
}

// Warn will log at the warn level using the logger associated with the
// context.
func Warn(ctx context.Context, msg string, fields ...logger.Field) {
	From(ctx).Warn(msg, fields...)
}

// Error will log at the error level using the logger associated with the
// context.
func Error(ctx context.Context, msg string, fields ...logger.Field) {
	From(ctx).Error(msg, fields...)
}

// Flush will flush any underlying log entries.
// applications should ensure that Flush is called before exiting.
func Flush(ctx context.Context) {
	From(ctx).Flush()
}
