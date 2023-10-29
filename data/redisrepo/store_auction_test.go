package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreAuction(t *testing.T) {

	matchOne := models.OrderFrag{
		OrderId: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Size:    decimal.NewFromFloat(200.0),
	}

	matchTwo := models.OrderFrag{
		OrderId: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Size:    decimal.NewFromFloat(300.0),
	}

	auctionID := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")
	auction := []models.OrderFrag{matchOne, matchTwo}

	auctionJson, _ := models.MarshalOrderFrags(auction)

	t.Run("should store auction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreateAuctionKey(auctionID), auctionJson).SetVal(1)

		err := repo.StoreAuction(ctx, auctionID, auction)
		assert.NoError(t, err)
	})

	t.Run("should return `ErrUnexpectedError` in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectSAdd(CreateAuctionKey(auctionID), auctionJson).SetErr(assert.AnError)

		err := repo.StoreAuction(ctx, auctionID, auction)
		assert.Equal(t, models.ErrUnexpectedError, err)
	})

}
