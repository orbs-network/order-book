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

	t.Run("QUOTE should return error zero amount in ", func(t *testing.T) {
		store := mocks.MockOrderBookStore{
			Error: models.ErrAuctionInvalid,
			Sets:  make(map[string]map[string]struct{}),
		}
		svc, _ := service.New(&store)
		_, err := svc.GetQuote(ctx, symbol, models.BUY, decimal.Zero)
		assert.Error(t, err, models.ErrInAmount)

	})

	t.Run("QUOTE HappyPath", func(t *testing.T) {

		mock := mocks.CreateAuctionMock()
		svc, _ := service.New(mock)
		ctx := context.Background()

		inAmount := decimal.NewFromInt(1000)
		outAmount := decimal.NewFromInt(1)
		res, err := svc.GetQuote(ctx, symbol, models.BUY, inAmount)
		assert.Equal(t, res.Size, outAmount)
		assert.NoError(t, err)
		//for _, frag := range res.OrderFrags {
		//order, err := svc.GetOrderById(ctx, frag.OrderId)
		//assert.NoError(t, err)
		//}
	})
}

// func TestService_RevertAuction(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("RevertAuction HappyPath", func(t *testing.T) {
// 		auctionId, _ := uuid.NewUUID()
// 		mock := mocks.CreateAuctionMock()
// 		svc, _ := service.New(mock)
// 		res, err := svc.ConfirmAuction(ctx, auctionId)
// 		assert.NoError(t, err)

// 		// all orders have pending size - no greater than the order itself
// 		for _, order := range res.Orders {
// 			assert.True(t, order.Size.GreaterThanOrEqual(order.SizePending))
// 			assert.True(t, order.SizePending.GreaterThan(decimal.Zero))
// 		}

// 		err = svc.RevertAuction(ctx, auctionId)
// 		assert.NoError(t, err)

// 		// all orders should not have pending size
// 		for _, order := range res.Orders {
// 			updatedOrder, err := svc.GetStore().FindOrderById(ctx, order.Id, false)
// 			assert.NoError(t, err)
// 			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
// 		}
// 	})
// }
// func TestService_AuctionMined(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("AuctionMined HappyPath", func(t *testing.T) {
// 		auctionId, _ := uuid.NewUUID()
// 		// creates the auction
// 		mock := mocks.CreateAuctionMock()
// 		svc, _ := service.New(mock)
// 		// confirm
// 		res, err := svc.ConfirmAuction(ctx, auctionId)
// 		assert.NoError(t, err)
// 		// mined
// 		err = svc.AuctionMined(ctx, auctionId)
// 		assert.NoError(t, err)

// 		// all orders should not have pending size
// 		for _, order := range res.Orders {
// 			updatedOrder, err := svc.GetStore().FindOrderById(ctx, order.Id, false)
// 			assert.NoError(t, err)
// 			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
// 		}
// 	})

// }
