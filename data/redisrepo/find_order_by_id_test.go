package redisrepo

import (
	"context"
	"errors"
	"fmt"
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

	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000008")
	clientOId := uuid.MustParse("00000000-0000-0000-0000-000000000009")
	symbol, _ := models.StrToSymbol("USDC-ETH")
	price := decimal.NewFromFloat(10000.55)
	zero := decimal.NewFromFloat(0)
	order := models.Order{
		Id:          orderID,
		ClientOId:   clientOId,
		UserId:      uuid.New(),
		Price:       price,
		Size:        sizeDec,
		SizePending: zero,
		SizeFilled:  zero,
		Symbol:      symbol,
		Side:        models.BUY,
		Status:      models.STATUS_OPEN,
	}

	t.Run("order found by orderID", func(t *testing.T) {
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetVal(order.OrderToMap())

		foundOrder, err := repo.FindOrderById(ctx, orderID, false)
		assert.NoError(t, err)
		fmt.Println(foundOrder)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("order found by clientOId", func(t *testing.T) {
		mock.ExpectGet(CreateClientOIDKey(clientOId)).SetVal(orderID.String())
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetVal(order.OrderToMap())

		foundOrder, err := repo.FindOrderById(ctx, clientOId, true)
		assert.NoError(t, err)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("order not found by orderId", func(t *testing.T) {
		nonExistentOrderID := uuid.New()
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetVal(map[string]string{})

		foundOrder, err := repo.FindOrderById(ctx, nonExistentOrderID, false)
		assert.Error(t, err, models.ErrOrderNotFound.Error())
		assert.Nil(t, foundOrder)
	})

	t.Run("order not found by clientOId", func(t *testing.T) {
		nonExistentClientOId := uuid.New()
		mock.ExpectGet(CreateClientOIDKey(nonExistentClientOId)).RedisNil()

		foundOrder, err := repo.FindOrderById(ctx, nonExistentClientOId, true)
		assert.Error(t, err, models.ErrOrderNotFound.Error())
		assert.Nil(t, foundOrder)
	})

	t.Run("invalid order ID format retrieved by clientOId", func(t *testing.T) {
		invalidOrderID := "invalid-order-id"
		mock.ExpectGet(CreateClientOIDKey(clientOId)).SetVal(invalidOrderID)

		foundOrder, err := repo.FindOrderById(ctx, clientOId, true)

		assert.Error(t, err, fmt.Sprintf("invalid order ID format retrieved by clientOId: %s", err))
		assert.Nil(t, foundOrder)
	})

	t.Run("unexpected error", func(t *testing.T) {
		var redisErr = errors.New("unexpected error")
		mock.ExpectHGetAll(CreateOrderIDKey(orderID)).SetErr(redisErr)

		foundOrder, err := repo.FindOrderById(ctx, orderID, false)
		assert.Error(t, err, redisErr.Error())
		assert.Nil(t, foundOrder)
	})
}
