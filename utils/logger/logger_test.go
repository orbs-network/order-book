package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func newTestLogger() (Logger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)

	return Logger{
		zap:           observedLogger,
		defaultFields: make(map[string]Field),
	}, observedLogs
}

func TestLogger(t *testing.T) {
	t.Run("with logger adds new fields", func(t *testing.T) {
		l, _ := newTestLogger()
		assert.Len(t, l.defaultFields, 0)

		f1 := String("key1", "val1")
		l = WithFields(l, f1)
		assert.Len(t, l.defaultFields, 1)

		f2 := String("key2", "val2")
		l = WithFields(l, f2)
		assert.Len(t, l.defaultFields, 2)

		// When we update an existing key with a new value the number
		// of new fields is not increased.
		f3 := String("key1", "val3")
		l = WithFields(l, f3)
		assert.Len(t, l.defaultFields, 2)

		assertFieldIn(t, l.defaultFields, f2)
		assertFieldIn(t, l.defaultFields, f3)
	})

	t.Run("with logger does not modify existing loggers default fields", func(t *testing.T) {
		l1, _ := newTestLogger()
		assert.Len(t, l1.defaultFields, 0)

		addedField := String("key1", "val1")
		l2 := WithFields(l1, addedField)

		assertFieldIn(t, l2.defaultFields, addedField)
		assertFieldNotIn(t, l1.defaultFields, addedField)
	})
}
