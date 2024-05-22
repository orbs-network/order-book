package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
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

	user := models.User{
		Id:     userId,
		PubKey: userPubKey,
		Type:   models.MARKET_MAKER,
	}

	input := service.CreateOrderInput{
		UserId:        userId,
		Price:         price,
		Symbol:        symbol,
		Size:          size,
		Side:          models.SELL,
		ClientOrderID: orderId,
		Eip712Sig:     "mock-sig",
		AbiFragment:   mocks.AbiFragment,
	}

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
		assert.Equal(t, newOrder.Signature, models.Signature{Eip712Sig: "mock-sig", AbiFragment: mocks.AbiFragment})
		assert.Equal(t, newOrder.Side, models.SELL)
	})

	t.Run("existing order with different userId - should return `ErrClashingOrderId` error", func(t *testing.T) {
		input := service.CreateOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
			Eip712Sig:     "mock-sig",
			AbiFragment:   mocks.AbiFragment,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{UserId: uuid.MustParse("b577273e-12de-4acc-a4f8-de7fb5b86e37")}}, mockBcClient)

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrClashingOrderId)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("existing order with same clientOrderId - should return `ErrClashingClientOrderId` error", func(t *testing.T) {
		input := service.CreateOrderInput{
			UserId:        userId,
			Price:         price,
			Symbol:        symbol,
			Size:          size,
			Side:          models.SELL,
			ClientOrderID: orderId,
			Eip712Sig:     "mock-sig",
			AbiFragment:   mocks.AbiFragment,
		}

		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, Order: &models.Order{ClientOId: orderId, UserId: userId}}, mockBcClient)

		order, err := svc.CreateOrder(ctx, input)

		assert.ErrorIs(t, err, models.ErrClashingClientOrderId)
		assert.Equal(t, models.Order{}, order)
	})

	t.Run("cross trade - bid higher than minAsk or ask lower than maxBid should fail", func(t *testing.T) {
		input := service.CreateOrderInput{
			UserId:      userId,
			Price:       price,
			Symbol:      symbol,
			Size:        size,
			Side:        models.BUY,
			Eip712Sig:   "mock-sig",
			AbiFragment: mocks.AbiFragment,
		}

		// create Bid
		size = decimal.NewFromInt(100)

		depth := models.MarketDepth{
			Asks: [][]decimal.Decimal{{decimal.NewFromFloat(0.11), size}},
			Bids: [][]decimal.Decimal{{decimal.NewFromFloat(0.9), size}},
		}
		svc, _ := service.New(&mocks.MockOrderBookStore{User: &user, MarketDepth: depth}, mockBcClient)

		// create Lower Ask should fail
		input.Side = models.SELL

		// above should work
		input.Price = decimal.NewFromFloat(0.95)
		_, err := svc.CreateOrder(ctx, input)
		assert.NoError(t, err)

		// below should fail
		input.Price = decimal.NewFromFloat(0.8)
		_, err = svc.CreateOrder(ctx, input)
		assert.ErrorIs(t, err, models.ErrCrossTrade)

		// create Higher Bid should fail
		input.Side = models.BUY

		// below should work
		input.Price = decimal.NewFromFloat(0.10)
		_, err = svc.CreateOrder(ctx, input)
		assert.NoError(t, err)

		// above should fail
		input.Price = decimal.NewFromFloat(0.12)
		_, err = svc.CreateOrder(ctx, input)
		assert.ErrorIs(t, err, models.ErrCrossTrade)

	})
}
