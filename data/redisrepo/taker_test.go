package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetSwap(t *testing.T) {

	swapJson := []string{
		`{"frags":[{"orderId":"550e8400-e29b-41d4-a716-446655440000","inSize":"10.5"},{"orderId":"550e8400-e29b-41d4-a716-446655440001","inSize":"5.3"}],"created":"2024-02-20T15:36:57.316643+02:00"}`,
		`{"frags":[{"orderId":"550e8400-e29b-41d4-a716-446655440002","inSize":"7.8"}],"created":"2024-02-20T15:36:57.316643+02:00"}`,
	}

	swapId := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")
	t.Run("should get swap", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectGet(CreateOpenSwapKey(swapId)).SetVal(swapJson[0])

		swap, err := repo.GetSwap(ctx, swapId)
		assert.NoError(t, err)
		assert.Len(t, swap.Frags, 2, "Should have 2 orders in the swap")
	})

	t.Run("should return `ErrUnexpectedError` in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectSMembers(CreateOpenSwapKey(swapId)).SetErr(assert.AnError)

		_, err := repo.GetSwap(ctx, swapId)
		assert.Equal(t, models.ErrUnexpectedError, err)
	})
}

func TestRedisRepository_RemoveSwap(t *testing.T) {
	swapId := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")

	t.Run("should remove swap", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectDel(CreateOpenSwapKey(swapId)).SetVal(1)

		err := repo.RemoveSwap(ctx, swapId)

		assert.NoError(t, err)
	})

	t.Run("should return error in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectDel(CreateOpenSwapKey(swapId)).SetErr(assert.AnError)

		err := repo.RemoveSwap(ctx, swapId)

		assert.Equal(t, assert.AnError, err)
	})
}
