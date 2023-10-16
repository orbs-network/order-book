package redisrepo

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redismock/v9"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

var symbol, _ = models.StrToSymbol("USDC-ETH")
var price = decimal.NewFromFloat(10.0)

var buyOrder = models.Order{
	Id:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	Price:  price,
	Side:   models.BUY,
	Symbol: symbol,
}

var sellOrder = models.Order{
	Id:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	Price:  price,
	Side:   models.SELL,
	Symbol: symbol,
}

var float64Price, _ = price.Float64()

func TestRedisRepository_GetBestPriceFor(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()

	repo := &redisRepository{
		client: db,
	}

	t.Run("BUY side - existing orders - order should be returned", func(t *testing.T) {

		buyPricesKey := CreateBuySidePricesKey(buyOrder.Symbol)

		mock.ExpectZRevRangeWithScores(buyPricesKey, 0, 0).SetVal([]redis.Z{{
			Score:  float64Price,
			Member: buyOrder.Id,
		}})

		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.BUY)

		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, buyOrder.Price, bestPrice, "prices should match")
	})

	t.Run("BUY side - no orders - zero should be returned with error", func(t *testing.T) {
		buyPricesKey := CreateBuySidePricesKey(buyOrder.Symbol)

		mock.ExpectZRevRangeWithScores(buyPricesKey, 0, 0).SetVal([]redis.Z{})

		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.BUY)

		assert.Error(t, err, models.ErrOrderNotFound, "error should be ErrOrderNotFound")
		assert.Equal(t, decimal.Zero, bestPrice, "should be zero")
	})

	t.Run("SELL side - existing orders - order should be returned", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)

		mock.ExpectZRangeWithScores(sellPricesKey, 0, 0).SetVal([]redis.Z{{
			Score:  float64Price,
			Member: sellOrder.Id,
		}})

		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, sellOrder.Price, bestPrice, "prices should match")
	})

	t.Run("SELL side - no orders - zero should be returned with error", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)

		mock.ExpectZRangeWithScores(sellPricesKey, 0, 0).SetVal([]redis.Z{})
		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.Error(t, err, models.ErrOrderNotFound, "error should be ErrOrderNotFound")
		assert.Equal(t, decimal.Zero, bestPrice, "should be zero")
	})

	t.Run("error with redis query - zero should be returned with error", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)

		someError := fmt.Errorf("something unexpected happened")

		mock.ExpectZRangeWithScores(sellPricesKey, 0, 0).SetErr(someError)
		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.Error(t, err, someError, "should be someError")
		assert.Equal(t, decimal.Zero, bestPrice, "should be zero")
	})

	t.Run("invalid side - zero should be returned with error", func(t *testing.T) {

		bestPrice, err := repo.GetBestPriceFor(ctx, symbol, models.Side("invalid"))

		assert.Error(t, err, "error should be returned")
		assert.Equal(t, decimal.Zero, bestPrice, "should be zero")
	})

}
