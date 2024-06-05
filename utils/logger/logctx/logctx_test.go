package logctx_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func WithObserver(ctx context.Context, opts ...zap.Option) (*observer.ObservedLogs, context.Context) {
	return WithObserverLevelEnabler(ctx, zap.DebugLevel, opts...)
}

func WithObserverLevelEnabler(ctx context.Context, enab zapcore.LevelEnabler, opts ...zap.Option) (*observer.ObservedLogs, context.Context) {
	zl, o := observer.New(enab)
	l := logger.New(logger.WithZapLogger(zap.New(zl, opts...)))
	return o, logctx.With(ctx, l)
}

func TestFrom(t *testing.T) {
	t.Run("should return a non-nil logger given there isn't a logger in the context", func(t *testing.T) {
		assert.NotNil(t, logctx.From(context.Background()))
	})

	t.Run("should return the zap logger we put into it", func(t *testing.T) {
		var (
			l   = logger.New(logger.WithZapLogger(zap.NewNop()))
			ctx = logctx.With(context.Background(), l)
		)
		assert.Equal(t, l, logctx.From(ctx))
	})

	t.Run("able to differentiate multiple contexts", func(t *testing.T) {
		const firstMessage = "some-first-message"
		const secondMessage = "some-second-message"

		o1, ctx1 := WithObserver(context.Background())
		o2, ctx2 := WithObserver(context.Background())

		ctx1 = logctx.WithFields(ctx1, logger.Bool("first_bool", true), logger.String("first_string", "A"))
		ctx2 = logctx.WithFields(ctx2, logger.Bool("second_bool", false), logger.String("second_string", "B"))

		logctx.Info(ctx1, firstMessage)
		logctx.Info(ctx2, secondMessage)

		m1 := o1.FilterMessage(firstMessage)
		require.NotZero(t, m1.Len())
		e1 := m1.TakeAll()[0]
		assert.Equal(t, zapcore.InfoLevel, e1.Level)
		assert.ElementsMatch(t, []zap.Field{zap.Bool("first_bool", true), zap.String("first_string", "A"), zap.String("service", "orderbook"), zap.String("commit_sha", "")}, e1.Context)

		m2 := o2.FilterMessage(secondMessage)
		require.NotZero(t, m2.Len())
		e2 := m2.TakeAll()[0]
		assert.Equal(t, zapcore.InfoLevel, e2.Level)
		assert.ElementsMatch(t, []zap.Field{zap.Bool("second_bool", false), zap.String("second_string", "B"), zap.String("service", "orderbook"), zap.String("commit_sha", "")}, e2.Context)
	})
}

func TestWith(t *testing.T) {
	t.Run("should return the zap logger we put into it", func(t *testing.T) {
		var (
			l   = logger.New(logger.WithZapLogger(zap.NewNop()))
			ctx = logctx.With(context.Background(), l)
		)
		assert.Equal(t, l, logctx.From(ctx))
	})
}

func TestWithFields(t *testing.T) {
	t.Run("should return the original context if there's nothing to do", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, ctx, logctx.WithFields(ctx))
	})
}

func TestDebug(t *testing.T) {
	t.Run("should call the zap logger associated with the context at the debug level", func(t *testing.T) {
		const someMessage = "some-message"

		o, ctx := WithObserver(context.Background())

		logctx.Debug(ctx, someMessage)

		l := o.FilterMessage(someMessage)
		require.NotZero(t, l.Len())
		assert.Equal(t, zapcore.DebugLevel, l.TakeAll()[0].Level)
	})
}

func TestInfo(t *testing.T) {
	t.Run("should call the zap logger associated with the context at the info level", func(t *testing.T) {
		const someMessage = "some-message"

		o, ctx := WithObserver(context.Background())

		logctx.Info(ctx, someMessage)

		l := o.FilterMessage(someMessage)
		require.NotZero(t, l.Len())
		assert.Equal(t, zapcore.InfoLevel, l.TakeAll()[0].Level)
	})
}

func TestWarn(t *testing.T) {
	t.Run("should call the zap logger associated with the context at the warn level", func(t *testing.T) {
		const someMessage = "some-message"

		o, ctx := WithObserver(context.Background())

		logctx.Warn(ctx, someMessage)

		l := o.FilterMessage(someMessage)
		require.NotZero(t, l.Len())
		assert.Equal(t, zapcore.WarnLevel, l.TakeAll()[0].Level)
	})
}

func TestError(t *testing.T) {
	t.Run("should call the zap logger associated with the context at the error level", func(t *testing.T) {
		const someMessage = "some-message"

		o, ctx := WithObserver(context.Background())

		logctx.Error(ctx, someMessage)

		l := o.FilterMessage(someMessage)
		require.NotZero(t, l.Len())
		assert.Equal(t, zapcore.ErrorLevel, l.TakeAll()[0].Level)
	})
}
