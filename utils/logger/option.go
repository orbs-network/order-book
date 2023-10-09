package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	zapLogger  *zap.Logger
	callerSkip int
}

// Option allows configuration of the logger instance.
type Option func(c *config)

func getConfig(opts []Option) config {
	cfg := config{
		zapLogger:  nil,
		callerSkip: 2,
	}

	for _, o := range opts {
		o(&cfg)
	}

	zl := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "@timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	), zap.AddCaller(), zap.AddCallerSkip(cfg.callerSkip))

	// Allow overriding the entire zap logger (for testing mainly).
	if cfg.zapLogger == nil {
		cfg.zapLogger = zl
	}

	return cfg
}

// WithZapLogger allows us to pass in a zap logger.
func WithZapLogger(l *zap.Logger) Option {
	return func(c *config) { c.zapLogger = l }
}

// WithCallerSkip allows us to skip additional levels when logging
// so that the caller is identified as the place in the code that calls the
// logger, and not the logger abstraction.
//
// For each level of abstraction over this library, add another skip level.
//
// default = 2 (skip this library)
func WithCallerSkip(skip int) Option {
	return func(c *config) { c.callerSkip = skip }
}
