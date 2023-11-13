package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_Utils_AddVal2Set(t *testing.T) {

	t.Run("should return without error when a new element is added that does not already exist", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectSAdd("key", "val").SetVal(1)

		err := AddVal2Set(context.Background(), client, "key", "val")
		assert.NoError(t, err)
	})

	t.Run("should return `ErrValAlreadyInSet` error when element already exists", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectSAdd("key", "val").SetVal(0)

		err := AddVal2Set(context.Background(), client, "key", "val")
		assert.ErrorIs(t, err, models.ErrValAlreadyInSet)
	})

	t.Run("should return error when unexpected Redis error", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectSAdd("key", "val").SetErr(assert.AnError)

		err := AddVal2Set(context.Background(), client, "key", "val")
		assert.ErrorIs(t, err, assert.AnError)
	})
}
