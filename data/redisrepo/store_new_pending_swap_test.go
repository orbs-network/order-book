package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo_StoreNewPendingSwap(t *testing.T) {
	ctx := context.Background()

	t.Run("should return `ErrNotFound` error if existing swap not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectGet(CreateOpenSwapKey(mocks.SwapTx.SwapId)).SetVal("")

		err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("should return error if redis fails to get swap", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectGet(CreateOpenSwapKey(mocks.SwapTx.SwapId)).SetErr(assert.AnError)

		err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

		assert.ErrorContains(t, err, "failed to get swap unexpectedly")
	})

}
