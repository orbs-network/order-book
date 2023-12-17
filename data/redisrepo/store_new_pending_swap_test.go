package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo_StoreNewPendingSwap(t *testing.T) {
	ctx := context.Background()

	t.Run("should add new pending swap to list", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreatePendingSwapTxsKey(), mocks.Pending.PendingToMap()).SetVal(1)

		err := repo.StoreNewPendingSwap(ctx, mocks.Pending)

		assert.NoError(t, err)
	})

	t.Run("should return error if redis fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreatePendingSwapTxsKey(), mocks.Pending.PendingToMap()).SetErr(assert.AnError)

		err := repo.StoreNewPendingSwap(ctx, mocks.Pending)

		assert.ErrorContains(t, err, "failed to store pending swap tx")
	})
}
