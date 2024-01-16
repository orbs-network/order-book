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

func TestService_CreateOrder(t *testing.T) {

	ctx := mocks.AddUserToCtx(nil)
	mockBcClient := &mocks.MockBcClient{IsVerified: true}

	symbol, _ := models.StrToSymbol("MATIC-USDC")
	userId := uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37")
	userPubKey := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
	price := decimal.NewFromFloat(10.0)
	orderId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")
	size := decimal.NewFromFloat(1000.00)

	eip712Domain := map[string]interface{}{}
	eip712MsgTypes := map[string]interface{}{}
	eip712Msg := map[string]interface{}{}

	user := models.User{
		Id:     userId,
		PubKey: userPubKey,
		Type:   models.MARKET_MAKER,
	}

	input := service.CreateOrderInput{
		UserId:         userId,
		Price:          price,
		Symbol:         symbol,
		Size:           size,
		Side:           models.SELL,
		ClientOrderID:  orderId,
		Eip712Sig:      "mock-sig",
		Eip712Domain:   &eip712Domain,
		Eip712MsgTypes: &eip712MsgTypes,
		Eip712Msg:      &eip712Msg,
	}

	t.Run("no user in context - should return error", func(t *testing.T) {

		ctxWithoutUser := context.Background()

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user}, mockBcClient)
		order, err := svc.CreateOrder(ctxWithoutUser, input)

		assert.ErrorContains(t, err, "user should be in context")
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("signature verification error - should return `ErrSignatureVerificationError` error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user}, &mocks.MockBcClient{Error: assert.AnError, IsVerified: false})

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrSignatureVerificationError)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("signature verification failed - should return `ErrSignatureVerificationFailed` error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user}, &mocks.MockBcClient{IsVerified: false})

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrSignatureVerificationFailed)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("unexpected error from store - should return error", func(t *testing.T) {

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Error: assert.AnError}, mockBcClient)

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorContains(t, err, "unexpected error when finding order by clientOrderId")
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("no previous order - should create new order", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: nil}, mockBcClient)

		newOrder, err := svc.CreateOrder(ctx, input)

		assert.NoError(t, err)
		// TODO: I am not asserting against the full order as timestamp is always different
		assert.NotEqual(t, newOrder.Id, orderId)
		assert.Equal(t, newOrder.ClientOId, orderId)
		assert.Equal(t, newOrder.UserId, user.Id)
		assert.Equal(t, newOrder.Price, price)
		assert.Equal(t, newOrder.Symbol, symbol)
		assert.Equal(t, newOrder.Size, size)
		assert.Equal(t, newOrder.Signature, models.Signature{Eip712Sig: "mock-sig", Eip712Domain: eip712Domain, Eip712MsgTypes: eip712MsgTypes, Eip712Msg: eip712Msg})
		assert.Equal(t, newOrder.Side, models.SELL)
	})

	t.Run("existing order with different userId - should return `ErrClashingOrderId` error", func(t *testing.T) {
		input := service.CreateOrderInput{
			UserId:         userId,
			Price:          price,
			Symbol:         symbol,
			Size:           size,
			Side:           models.SELL,
			ClientOrderID:  orderId,
			Eip712Sig:      "mock-sig",
			Eip712Domain:   &eip712Domain,
			Eip712MsgTypes: &eip712MsgTypes,
			Eip712Msg:      &eip712Msg,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{UserId: uuid.MustParse("b577273e-12de-4acc-a4f8-de7fb5b86e37")}}, mockBcClient)

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrClashingOrderId)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("existing order with same clientOrderId - should return `ErrClashingClientOrderId` error", func(t *testing.T) {
		input := service.CreateOrderInput{
			UserId:         userId,
			Price:          price,
			Symbol:         symbol,
			Size:           size,
			Side:           models.SELL,
			ClientOrderID:  orderId,
			Eip712Sig:      "mock-sig",
			Eip712Domain:   &eip712Domain,
			Eip712MsgTypes: &eip712MsgTypes,
			Eip712Msg:      &eip712Msg,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{ClientOId: orderId, UserId: userId}}, mockBcClient)

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrClashingClientOrderId)
		assert.Equal(t, models.Order{}, order)
	})
}
