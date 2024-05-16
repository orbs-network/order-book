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

const symbol = "ETH-USD"

// /////////////////////////////////////////////////////////////
func TestTaker_Quote(t *testing.T) {
	ctx := context.Background()
	evmClient := &service.EvmClient{}

	t.Run("QUOTE should return error zero amount in", func(t *testing.T) {
		store := mocks.MockOrderBookStore{
			Error: models.ErrSwapInvalid,
			Sets:  make(map[string]map[string]struct{}),
		}
		// buy
		svc, _ := service.New(&store, evmClient)
		res, err := svc.GetQuote(ctx, symbol, models.BUY, decimal.Zero, nil, "0xTOKEN")
		assert.Equal(t, res, models.QuoteRes{})
		assert.Error(t, err, models.ErrInAmount)
		// sell
		svc, _ = service.New(&store, evmClient)
		res, err = svc.GetQuote(ctx, symbol, models.SELL, decimal.Zero, nil, "0xTOKEN")
		assert.Equal(t, res, models.QuoteRes{})
		assert.Error(t, err, models.ErrInAmount)

	})

	// t.Run("QUOTE HappyPath buy", func(t *testing.T) {
	// 	mock := mocks.CreateSwapMock()
	// 	svc, _ := service.New(mock, evmClient)

	// 	inAmount := decimal.NewFromInt(1000)
	// 	outAmount := decimal.NewFromInt(1)
	// 	res, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount, nil, 18, 6)
	// 	assert.True(t, res.Size.Equals(outAmount))
	// 	assert.NoError(t, err)
	// })

	// t.Run("QUOTE HappyPath Sell", func(t *testing.T) {
	// 	mock := mocks.CreateSwapMock()
	// 	svc, _ := service.New(mock, evmClient)

	// 	inAmount := decimal.NewFromInt(1)
	// 	outAmount := decimal.NewFromInt(900)
	// 	res, err := svc.GetQuote(ctx, symbol, models.SELL, inAmount, nil, 6, 18)
	// 	assert.True(t, res.Size.Equals(outAmount))
	// 	assert.NoError(t, err)
	// })
}

// func TestTaker_BeginSwap(t *testing.T) {
// 	ctx := context.Background()
// 	evmClient := &service.EvmClient{}

// 	t.Run("BeginSwap Should return the same as quote, second quote returns diff amount", func(t *testing.T) {
// 		mock := mocks.CreateSwapMock()
// 		svc, _ := service.New(mock, evmClient)

// 		// get quote does not lock liquidity
// 		inAmount := decimal.NewFromInt(1000)
// 		outAmount := decimal.NewFromInt(1)
// 		oaRes, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount, nil, 18, 6)
// 		assert.True(t, oaRes.Size.Equals(outAmount))
// 		assert.NoError(t, err)

// 		swapRes, err := svc.BeginSwap(ctx, oaRes)
// 		assert.NoError(t, err)
// 		assert.Greater(t, len(swapRes.Orders), 0)
// 		assert.Greater(t, len(swapRes.Fragments), 0)
// 		// amount out should be equal in quote and swap requests
// 		assert.True(t, oaRes.Size.Equals(swapRes.OutAmount))

// 		// second quote however should return different outAmount as first order has already been filled
// 		oaRes2, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount, nil, 6, 18)
// 		assert.NoError(t, err)
// 		assert.NotEqual(t, oaRes.Size, oaRes2.Size)
// 	})
// }

func TestTaker_SwapStarted(t *testing.T) {
	ctx := context.Background()
	evmClient := &service.EvmClient{}
	mock := mocks.CreateSwapMock()
	svc, _ := service.New(mock, evmClient)

	err := svc.SwapStarted(ctx, uuid.New(), "0x123334")
	assert.NoError(t, err)
}

// func TestService_AbortSwap(t *testing.T) {
// 	ctx := context.Background()
// 	evmClient := &service.EvmClient{}

// 	t.Run("AbortSwap HappyPath", func(t *testing.T) {
// 		mock := mocks.CreateSwapMock()
// 		svc, _ := service.New(mock, evmClient)

// 		inAmount := decimal.NewFromInt(1000)
// 		outAmount := decimal.NewFromInt(1)
// 		oaRes, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount, nil, 18, 6)
// 		assert.True(t, oaRes.Size.Equals(outAmount))
// 		assert.NoError(t, err)

// 		swapRes, err := svc.BeginSwap(ctx, oaRes)
// 		assert.NoError(t, err)
// 		assert.Greater(t, len(swapRes.Orders), 0)
// 		assert.Greater(t, len(swapRes.Fragments), 0)

// 		// all orders have pending size - no greater than the order itself
// 		for _, order := range swapRes.Orders {
// 			assert.True(t, order.Size.GreaterThanOrEqual(order.SizePending))
// 			assert.True(t, order.SizePending.GreaterThan(decimal.Zero))
// 		}

// 		err = svc.AbortSwap(ctx, swapRes.SwapId)
// 		assert.NoError(t, err)

// 		// all orders should not have pending size
// 		for _, order := range swapRes.Orders {
// 			updatedOrder, err := svc.GetOrderById(ctx, order.Id)
// 			assert.NoError(t, err)
// 			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
// 		}
// 	})
// }
