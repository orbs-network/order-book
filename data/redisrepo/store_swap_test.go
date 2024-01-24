package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreSwap(t *testing.T) {

	matchOne := models.OrderFrag{
		OrderId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		OutSize: decimal.NewFromFloat(200.0),
	}

	matchTwo := models.OrderFrag{
		OrderId: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		OutSize: decimal.NewFromFloat(300.0),
	}

	swapID := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")
	frags := []models.OrderFrag{matchOne, matchTwo}

	swapJson, _ := models.MarshalOrderFrags(frags)

	t.Run("should store swap", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreateSwapKey(swapID), swapJson).SetVal(1)

		err := repo.StoreSwap(ctx, swapID, frags)
		assert.NoError(t, err)
	})

	t.Run("should return `ErrUnexpectedError` in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreateSwapKey(swapID), swapJson).SetErr(assert.AnError)

		err := repo.StoreSwap(ctx, swapID, frags)
		assert.ErrorContains(t, err, "failed to store swap")
	})

}
