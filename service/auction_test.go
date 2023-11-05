package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

// /////////////////////////////////////////////////////////////////
// static iterator impl
// type iter struct {
// 	orders []models.Order
// 	index  int
// }

// func (i *iter) HasNext() bool {
// 	return i.index < (len(i.orders) - 1)
// }
// func (i *iter) Next(ctx context.Context) *models.Order {
// 	// increment index
// 	i.index = i.index + 1

// 	if i.index >= len(i.orders) {
// 		return nil
// 	}

// 	// get order
// 	return &i.orders[i.index]
// }

// ///////////////////////////////////////////////////////////////
func TestService_ConfirmAuction(t *testing.T) {
	ctx := context.Background()

	t.Run("ConfirmAuction error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{Error: assert.AnError})
		uuid, _ := uuid.NewUUID()
		_, err := svc.ConfirmAuction(ctx, uuid)
		assert.Error(t, err)

	})

	t.Run("ConfirmAuction HappyPath", func(t *testing.T) {
		uuid, _ := uuid.NewUUID()
		mock := mocks.CreateAuctionMock()
		svc, _ := service.New(mock)
		res, err := svc.ConfirmAuction(ctx, uuid)
		assert.NoError(t, err)
		for i := 1; i < len(res.Orders); i++ {
			order := res.Orders[i]
			frag := res.Fragments[i]
			assert.Equal(t, frag.OrderId.String(), order.Id.String())
			assert.Equal(t, frag.Size, order.SizePending)
		}

		// make sure last fraf size is not equal to size
		last := len(res.Orders) - 1
		order := res.Orders[last]
		frag := res.Fragments[last]
		assert.NotEqual(t, frag.Size, order.Size)
	})
}
func TestService_RevertAuction(t *testing.T) {
	ctx := context.Background()

	t.Run("RevertAuction HappyPath", func(t *testing.T) {
		auctionId, _ := uuid.NewUUID()
		mock := mocks.CreateAuctionMock()
		svc, _ := service.New(mock)
		res, err := svc.ConfirmAuction(ctx, auctionId)
		assert.NoError(t, err)

		// all orders have pending size - no greater than the order itself

		for _, order := range res.Orders {
			assert.True(t, order.Size.GreaterThanOrEqual(order.SizePending))
			assert.True(t, order.SizePending.GreaterThan(decimal.Zero))
		}

		err = svc.RevertAuction(ctx, auctionId)
		assert.NoError(t, err)

		// all orders should not have pending size
		for _, order := range res.Orders {
			updatedOrder, err := svc.GetStore().FindOrderById(ctx, order.Id, false)
			assert.NoError(t, err)
			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
		}
	})
}
func TestService_AuctionMined(t *testing.T) {
	ctx := context.Background()

	t.Run("AuctionMined HappyPath", func(t *testing.T) {
		auctionId, _ := uuid.NewUUID()
		// creates the auction
		mock := mocks.CreateAuctionMock()
		svc, _ := service.New(mock)
		// confirm
		res, err := svc.ConfirmAuction(ctx, auctionId)
		assert.NoError(t, err)
		// mined
		err = svc.AuctionMined(ctx, auctionId)
		assert.NoError(t, err)

		// all orders should not have pending size
		for _, order := range res.Orders {
			updatedOrder, err := svc.GetStore().FindOrderById(ctx, order.Id, false)
			assert.NoError(t, err)
			assert.True(t, updatedOrder.SizePending.Equal(decimal.Zero))
		}
	})

}
