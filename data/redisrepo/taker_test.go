package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetSwap(t *testing.T) {

	uuid1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	uuid2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	uuid3 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	amount1 := decimal.NewFromFloat(10.5)
	amount2 := decimal.NewFromFloat(5.3)
	amount3 := decimal.NewFromFloat(7.8)

	swapJson := []string{
		`[{"orderId":"550e8400-e29b-41d4-a716-446655440000","size":"10.5"},{"orderId":"550e8400-e29b-41d4-a716-446655440001","size":"5.3"}]`,
		`[{"orderId":"550e8400-e29b-41d4-a716-446655440002","size":"7.8"}]`,
	}

	swapId := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")
	t.Run("should get swap", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectLRange(CreateSwapKey(swapId), 0, -1).SetVal(swapJson)

		swap, err := repo.GetSwap(ctx, swapId)
		assert.NoError(t, err)
		assert.Len(t, swap, 3, "Should have 3 orders in the swap")
		assert.ElementsMatch(t, []models.OrderFrag{
			{OrderId: uuid1, Size: amount1},
			{OrderId: uuid2, Size: amount2},
			{OrderId: uuid3, Size: amount3},
		}, swap, "The swap contents do not match expected")
	})

	t.Run("should return `ErrUnexpectedError` in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectSMembers(CreateSwapKey(swapId)).SetErr(assert.AnError)

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

		mock.ExpectDel(CreateSwapKey(swapId)).SetVal(1)

		err := repo.RemoveSwap(ctx, swapId)

		assert.NoError(t, err)
	})

	t.Run("should return error in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectDel(CreateSwapKey(swapId)).SetErr(assert.AnError)

		err := repo.RemoveSwap(ctx, swapId)

		assert.Equal(t, assert.AnError, err)
	})
}
