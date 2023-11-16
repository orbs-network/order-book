package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestService_ProcessOrder(t *testing.T) {

	ctx := context.Background()

	symbol, _ := models.StrToSymbol("USDC-ETH")
	userId := uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37")
	userPubKey := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
	price := decimal.NewFromFloat(10.0)
	orderId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")
	size := decimal.NewFromFloat(1000.00)

	user := models.User{
		Id:     userId,
		PubKey: userPubKey,
		Type:   models.MARKET_MAKER,
	}

	t.Run("unexpected error from store - should return `ErrUnexpectedError` error", func(t *testing.T) {
		input := service.ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Error: assert.AnError})

		order, err := svc.ProcessOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrUnexpectedError)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("no previous order - should create new order", func(t *testing.T) {
		input := service.ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: nil})

		newOrder, err := svc.ProcessOrder(ctx, input)

		assert.NoError(t, err)
		// TODO: I am not asserting against the full order as timestamp is always different
		assert.NotEqual(t, newOrder.Id, orderId)
		assert.Equal(t, newOrder.ClientOId, orderId)
		assert.Equal(t, newOrder.UserId, user.Id)
		assert.Equal(t, newOrder.Price, price)
		assert.Equal(t, newOrder.Symbol, symbol)
		assert.Equal(t, newOrder.Size, size)
		assert.Equal(t, newOrder.Signature, "")
		assert.Equal(t, newOrder.Status, models.STATUS_OPEN)
		assert.Equal(t, newOrder.Side, models.SELL)
	})

	t.Run("existing order with different userId - should return `ErrClashingOrderId` error", func(t *testing.T) {
		input := service.ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{UserId: uuid.MustParse("b577273e-12de-4acc-a4f8-de7fb5b86e37")}})

		order, err := svc.ProcessOrder(ctx, input)

		assert.ErrorIs(t, err, service.ErrClashingOrderId)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("existing order with same clientOrderId - should return `ErrOrderAlreadyExists` error", func(t *testing.T) {
		input := service.ProcessOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{ClientOId: orderId, UserId: userId}})

		order, err := svc.ProcessOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrOrderAlreadyExists)
		assert.Equal(t, models.Order{}, order)
	})
}
