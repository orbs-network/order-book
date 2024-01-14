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

func TestGetMarketDepth(t *testing.T) {
	// Create a new Service with a mock orderBookStore
	ctx := context.Background()

	md := models.MarketDepth{
		Asks: [][]decimal.Decimal{
			{decimal.NewFromInt(1), decimal.NewFromInt(2)},
		},
		Bids: [][]decimal.Decimal{
			{decimal.NewFromInt(3), decimal.NewFromInt(4)},
		},
		Symbol: "MATIC-USDC",
		Time:   1634567890,
	}

	// Test case 1: Successful scenario
	t.Run("should successfully return market depth", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{
			MarketDepth: md,
		}, &mocks.MockBcClient{})

		marketDepth, err := svc.GetMarketDepth(ctx, "MATIC-USDC", 5)

		assert.NoError(t, err, "Get market depth should not return an error")
		assert.Equal(t, md, marketDepth, "Expected non-nil market depth")
	})

	// Test case 1: Error scenario
	t.Run("should return Error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{
			MarketDepth: md,
			Error:       assert.AnError,
		}, &mocks.MockBcClient{})

		marketDepth, err := svc.GetMarketDepth(ctx, "MATIC-USDC", 5)

		assert.Zero(t, len(marketDepth.Asks), "marketDepth.Asks should be empty")
		assert.Zero(t, len(marketDepth.Bids), "marketDepth.Bids should be empty")

		assert.Error(t, err, "Get market depth should not return an error")
	})
}
