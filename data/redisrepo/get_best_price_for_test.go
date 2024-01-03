package redisrepo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redismock/v9"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

var symbol, _ = models.StrToSymbol("MATIC-USDC")
var price = decimal.NewFromFloat(10.0)

var buyOrder = models.Order{
	Id:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	UserId: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	Price:  price,
	Size:   decimal.NewFromFloat(1212312.0),
	Symbol: symbol,
	Side:   models.BUY,
}

var sellOrder = models.Order{
	Id:     uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	UserId: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	Price:  price,
	Size:   decimal.NewFromFloat(1212312.0),
	Symbol: symbol,
	Side:   models.SELL,
}

func TestRedisRepository_GetBestPriceFor(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()

	repo := &redisRepository{
		client: db,
	}

	t.Run("BUY side - existing orders - order should be returned", func(t *testing.T) {

		buyPricesKey := CreateBuySidePricesKey(buyOrder.Symbol)
		mock.ExpectZRevRange(buyPricesKey, 0, 0).SetVal([]string{buyOrder.Id.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(buyOrder.Id)).SetVal(buyOrder.OrderToMap())

		order, err := repo.GetBestPriceFor(ctx, symbol, models.BUY)

		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, buyOrder.Price.String(), order.Price.String(), "prices should match")
	})

	t.Run("BUY side - no orders - zero should be returned with error", func(t *testing.T) {
		buyPricesKey := CreateBuySidePricesKey(buyOrder.Symbol)

		mock.ExpectZRevRange(buyPricesKey, 0, 0).SetVal([]string{})
		order, err := repo.GetBestPriceFor(ctx, symbol, models.BUY)

		assert.Error(t, err, models.ErrNotFound, "error should be ErrNotFound")
		assert.Equal(t, models.Order{}, order, "should be zero")
	})

	t.Run("SELL side - existing orders - order should be returned", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)
		mock.ExpectZRange(sellPricesKey, 0, 0).SetVal([]string{sellOrder.Id.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(sellOrder.Id)).SetVal(sellOrder.OrderToMap())

		order, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.NoError(t, err, "error should be nil")
		assert.Equal(t, sellOrder.Price.String(), order.Price.String(), "prices should match")
	})

	t.Run("SELL side - no orders - zero should be returned with error", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)

		mock.ExpectZRange(sellPricesKey, 0, 0).SetVal([]string{})
		order, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.Error(t, err, models.ErrNotFound, "error should be ErrNotFound")
		assert.Equal(t, models.Order{}, order, "should be zero")
	})

	t.Run("error with redis query - zero should be returned with error", func(t *testing.T) {
		sellPricesKey := CreateSellSidePricesKey(sellOrder.Symbol)

		mock.ExpectZRange(sellPricesKey, 0, 0).SetErr(assert.AnError)
		order, err := repo.GetBestPriceFor(ctx, symbol, models.SELL)

		assert.Error(t, err, assert.AnError, "should have errored")
		assert.Equal(t, models.Order{}, order, "should be zero value")
	})

	t.Run("invalid side - zero should be returned with error", func(t *testing.T) {

		order, err := repo.GetBestPriceFor(ctx, symbol, models.Side("invalid"))

		assert.ErrorIs(t, err, ErrInvalidOrderSide)
		assert.Equal(t, models.Order{}, order, "should be zero value")
	})

}
