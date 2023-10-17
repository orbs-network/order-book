package redisrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_FindOrderById(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()

	repo := &redisRepository{
		client: db,
	}

	sizeDec, _ := decimal.NewFromString("126")

	orderID := uuid.New()
	symbol, _ := models.StrToSymbol("USDC-ETH")
	price := decimal.NewFromFloat(10000.55)
	order := models.Order{
		Id:     orderID,
		UserId: uuid.New(),
		Price:  price,
		Size:   sizeDec,
		Symbol: symbol,
		Side:   models.BUY,
		Status: models.STATUS_OPEN,
	}

	t.Run("order found", func(t *testing.T) {
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetVal(order.OrderToMap())

		foundOrder, err := repo.FindOrderById(ctx, orderID)
		assert.NoError(t, err)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("order not found", func(t *testing.T) {
		nonExistentOrderID := uuid.New()
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetVal(map[string]string{})

		foundOrder, err := repo.FindOrderById(ctx, nonExistentOrderID)
		assert.Error(t, err, models.ErrOrderNotFound.Error())
		assert.Nil(t, foundOrder)
	})

	t.Run("unexpected error", func(t *testing.T) {
		var redisErr = errors.New("unexpected error")
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetErr(redisErr)

		foundOrder, err := repo.FindOrderById(ctx, orderID)
		assert.Error(t, err, redisErr.Error())
		assert.Nil(t, foundOrder)
	})
}
