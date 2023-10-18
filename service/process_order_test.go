package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestService_ProcessOrder(t *testing.T) {

	ctx := context.Background()

	symbol, _ := models.StrToSymbol("USDC-ETH")

	userId := uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37")
	price := decimal.NewFromFloat(10.0)

	orderId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")
	clientOrderId := uuid.MustParse("d577273e-12de-4acc-a4f8-de7fb5b86e37")

	order := models.Order{
		Id:            orderId,
		UserId:        userId,
		Price:         price,
		Symbol:        symbol,
		Size:          decimal.NewFromFloat(1.0),
		Signature:     "",
		Status:        models.STATUS_OPEN,
		Side:          models.SELL,
		ClientOrderID: clientOrderId,
	}

	t.Run("new order with clientOrderId", func(t *testing.T) {
		input := ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          decimal.NewFromFloat(1.0),
			Side:          models.SELL,
			ClientOrderID: &clientOrderId,
		}

		svc, _ := New(&mocks.MockOrderBookStore{Order: order})

		order, err := svc.ProcessOrder(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, userId, order.UserId)
		assert.Equal(t, price, order.Price)
		assert.Equal(t, symbol, order.Symbol)
		assert.Equal(t, decimal.NewFromFloat(1.0), order.Size)
		assert.Equal(t, models.SELL, order.Side)
		assert.NotEqual(t, uuid.Nil, order.Id)
		assert.Equal(t, models.STATUS_OPEN, order.Status)
		assert.Equal(t, clientOrderId, order.ClientOrderID)
	})

	t.Run("new order without ClientOrderID - orderId and clientOrderId should be set to the same ID", func(t *testing.T) {
		input := ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          decimal.NewFromFloat(1.0),
			Side:          models.SELL,
			ClientOrderID: nil,
		}

		orderWithSameIds := models.Order{
			Id:            orderId,
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          decimal.NewFromFloat(1.0),
			Signature:     "",
			Status:        models.STATUS_OPEN,
			Side:          models.SELL,
			ClientOrderID: orderId,
		}

		svc, _ := New(&mocks.MockOrderBookStore{Order: orderWithSameIds})

		order, err := svc.ProcessOrder(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, userId, order.UserId)
		assert.Equal(t, price, order.Price)
		assert.Equal(t, symbol, order.Symbol)
		assert.Equal(t, decimal.NewFromFloat(1.0), order.Size)
		assert.Equal(t, models.SELL, order.Side)
		assert.NotEqual(t, uuid.Nil, order.Id)
		assert.Equal(t, models.STATUS_OPEN, order.Status)
		assert.Equal(t, order.Id, order.ClientOrderID)
	})

	t.Run("process order with error from store", func(t *testing.T) {
		input := ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          decimal.NewFromFloat(1.0),
			Side:          models.SELL,
			ClientOrderID: nil,
		}

		svc, _ := New(&mocks.MockOrderBookStore{Error: assert.AnError})

		order, err := svc.ProcessOrder(ctx, input)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, models.Order{}, order)
	})
}

func TestGenerateClientOrderId(t *testing.T) {
	orderId := uuid.New()
	clientOrderId := uuid.New()

	// Case 1: clientOrderId is nil
	result := generateClientOrderId(nil, orderId)
	assert.Equal(t, orderId, result, "no clientOrderId passed so should be the same as order ID")

	// Case 2: clientOrderId is not nil
	result = generateClientOrderId(&clientOrderId, orderId)
	assert.Equal(t, clientOrderId, result, "clientOrderId passed so should use that")
}
