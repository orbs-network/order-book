// Package logger exposes common functionality for logging.
package logger

import (
	"errors"
	"os"
	"syscall"

	"go.uber.org/zap"
)

// Logger is an abstraction for logging in services.
// This hides any concrete implentation.e.g.
type Logger struct {
	zap           *zap.Logger
	defaultFields map[string]Field
}

// New creates a new production-ready logger instance.
func New(opts ...Option) Logger {
	cfg := getConfig(opts)

	commitSha := os.Getenv("COMMIT_SHA")

	defaultFields := map[string]Field{
		"service":    String("service", "orderbook"),
		"commit_sha": String("commit_sha", commitSha),
	}

	return Logger{
		zap:           cfg.zapLogger,
		defaultFields: defaultFields,
	}
}

// Debug logs at a debug level.
func (l Logger) Debug(msg string, fields ...Field) {
	l.log(l.zap.Debug, msg, fields)
}

// Info logs at an info level.
func (l Logger) Info(msg string, fields ...Field) {
	l.log(l.zap.Info, msg, fields)
}

// Warn logs at warn level.
func (l Logger) Warn(msg string, fields ...Field) {
	l.log(l.zap.Warn, msg, fields)
}

// Error logs at error level.
func (l Logger) Error(msg string, fields ...Field) {
	l.log(l.zap.Error, msg, fields)
}

func (l Logger) log(fn func(msg string, fields ...zap.Field), msg string, fields []Field) {
	fn(msg, mapToZap(l.defaultFields, fields)...)
}

// Flush will flush any underlying log entries.
// applications should ensure that Flush is called before exiting.
func (l Logger) Flush() {
	err := l.zap.Sync()

	// There is a known issue when writing to stdout/stderr as these
	// locations do not support flushing.  Flush works fine
	// when writing to files as is the case in most production applications.
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		l.Warn("error flushing logger", Error(err))
	}
}

// WithFields creates a new logger from the current logger with the fields
// automatically available in the case where a logger always needs a field.
//
// If a field is provided with an already existing key, the latest key provided
// will be given priority.
//
// To ensure that the new logger does not affect the old loggers default fields
// a true copy of the default fields map is taken. This is required as maps in Go
// are a pointer to the map location, so mutating the original map, will affect that
// stored on the original logger.
func WithFields(l Logger, fields ...Field) Logger {
	df := make(map[string]Field, len(l.defaultFields))
	for k, v := range l.defaultFields {
		df[k] = v
	}

	for i := range fields {
		df[fields[i].Key()] = fields[i]
	}

	return Logger{
		zap:           l.zap,
		defaultFields: df,
	}
}
