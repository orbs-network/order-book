package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetAuction(t *testing.T) {

	uuid1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	uuid2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	uuid3 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

	amount1 := decimal.NewFromFloat(10.5)
	amount2 := decimal.NewFromFloat(5.3)
	amount3 := decimal.NewFromFloat(7.8)

	auctionJson := []string{
		`[{"orderId":"550e8400-e29b-41d4-a716-446655440000","size":"10.5"},{"orderId":"550e8400-e29b-41d4-a716-446655440001","size":"5.3"}]`,
		`[{"orderId":"550e8400-e29b-41d4-a716-446655440002","size":"7.8"}]`,
	}

	auctionID := uuid.MustParse("a777273e-12de-4acc-a4f8-de7fb5b86e37")
	t.Run("should get auction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectLRange(CreateAuctionKey(auctionID), 0, -1).SetVal(auctionJson)

		auction, err := repo.GetAuction(ctx, auctionID)
		assert.NoError(t, err)
		assert.Len(t, auction, 3, "Should have 3 orders in the auction")
		assert.ElementsMatch(t, []models.OrderFrag{
			{OrderId: uuid1, Size: amount1},
			{OrderId: uuid2, Size: amount2},
			{OrderId: uuid3, Size: amount3},
		}, auction, "The auction contents do not match expected")
	})

	t.Run("should return `ErrUnexpectedError` in case of a Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectSMembers(CreateAuctionKey(auctionID)).SetErr(assert.AnError)

		_, err := repo.GetAuction(ctx, auctionID)
		assert.Equal(t, models.ErrUnexpectedError, err)
	})
}

// func TestRedisRepository_RemoveAuction(t *testing.T) {

// }
