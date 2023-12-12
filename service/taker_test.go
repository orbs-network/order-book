package service_test

import (
	"context"
	"testing"

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
	ethClient := &service.EthereumClient{}

	t.Run("QUOTE should return error zero amount in ", func(t *testing.T) {
		store := mocks.MockOrderBookStore{
			Error: models.ErrAuctionInvalid,
			Sets:  make(map[string]map[string]struct{}),
		}
		svc, _ := service.New(&store, ethClient)
		_, err := svc.GetQuote(ctx, symbol, models.BUY, decimal.Zero)
		assert.Error(t, err, models.ErrInAmount)

	})

	t.Run("QUOTE HappyPath", func(t *testing.T) {
		mock := mocks.CreateSwapMock()
		svc, _ := service.New(mock, ethClient)

		inAmount := decimal.NewFromInt(1000)
		outAmount := decimal.NewFromInt(1)
		res, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount)
		assert.True(t, res.Size.Equals(outAmount))
		assert.NoError(t, err)
	})
}
func TestTaker_BeginSwap(t *testing.T) {
	ctx := context.Background()
	ethClient := &service.EthereumClient{}

	t.Run("BeginSwap Should return the same as quote, second quote returns diff amount", func(t *testing.T) {
		mock := mocks.CreateSwapMock()
		svc, _ := service.New(mock, ethClient)

		// get quote does not lock liquidity
		inAmount := decimal.NewFromInt(1000)
		outAmount := decimal.NewFromInt(1)
		oaRes, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount)
		assert.True(t, oaRes.Size.Equals(outAmount))
		assert.NoError(t, err)

		swapRes, err := svc.BeginSwap(ctx, oaRes)
		assert.NoError(t, err)
		assert.Greater(t, len(swapRes.Orders), 0)
		assert.Greater(t, len(swapRes.Fragments), 0)
		// amount out should be equal in quote and swap requests
		assert.True(t, oaRes.Size.Equals(swapRes.OutAmount))

		// second quote however should return different amountOut as first order has already been filled
		oaRes2, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount)
		assert.NoError(t, err)
		assert.NotEqual(t, oaRes.Size, oaRes2.Size)
	})
}

func TestService_AbortSwap(t *testing.T) {
	ctx := context.Background()
	ethClient := &service.EthereumClient{}

	t.Run("AbortSwap HappyPath", func(t *testing.T) {
		mock := mocks.CreateSwapMock()
		svc, _ := service.New(mock, ethClient)

		inAmount := decimal.NewFromInt(1000)
		outAmount := decimal.NewFromInt(1)
		oaRes, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount)
		assert.True(t, oaRes.Size.Equals(outAmount))
		assert.NoError(t, err)

		swapRes, err := svc.BeginSwap(ctx, oaRes)
		assert.NoError(t, err)
		assert.Greater(t, len(swapRes.Orders), 0)
		assert.Greater(t, len(swapRes.Fragments), 0)

		// all orders have pending size - no greater than the order itself
		for _, order := range swapRes.Orders {
			assert.True(t, order.Size.GreaterThanOrEqual(order.SizePending))
			assert.True(t, order.SizePending.GreaterThan(decimal.Zero))
		}

		err = svc.AbortSwap(ctx, swapRes.SwapId)
		assert.NoError(t, err)

		// all orders should not have pending size
		for _, order := range swapRes.Orders {
			updatedOrder, err := svc.GetOrderById(ctx, order.Id)
			assert.NoError(t, err)
			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
		}
	})
}
