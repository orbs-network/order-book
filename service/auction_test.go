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

func newOrder(price, size int64) *models.Order {
	oid, _ := uuid.NewUUID()
	return &models.Order{
		Id:     oid,
		Price:  decimal.NewFromInt(price),
		Size:   decimal.NewFromInt(size),
		Status: models.STATUS_OPEN,
	}
}

func newAsks() []*models.Order {
	return []*models.Order{
		newOrder(1000, 1),
		newOrder(1001, 2),
		newOrder(1002, 3),
	}
}
func newBids() []*models.Order {
	return []*models.Order{
		newOrder(900, 1),
		newOrder(800, 2),
		newOrder(700, 3),
	}
}

func newFrags(orders []*models.Order) []models.OrderFrag {
	frags := []models.OrderFrag{}
	// create frag of all input orders except last one, which is only half filled
	for i, order := range orders {
		sz := order.Size
		// last element make half size
		if i == len(orders)-1 {
			sz = sz.Div(decimal.NewFromInt(2))
		}
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Size: sz})
	}
	return frags
}

func createAuctionMock() *mocks.MockOrderBookStore {
	res := mocks.MockOrderBookStore{Error: nil}
	res.Asks = newAsks()
	res.Bids = newBids()
	res.Frags = newFrags(res.Asks)
	return &res
}

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
		mock := createAuctionMock()
		svc, _ := service.New(mock)
		res, err := svc.ConfirmAuction(ctx, uuid)
		assert.NoError(t, err)
		for i := 1; i < len(res.Orders); i++ {
			order := res.Orders[i]
			frag := res.Fragments[i]
			assert.Equal(t, frag.OrderId.String(), order.Id.String())
			assert.Equal(t, frag.Size, order.SizePending)
		}
	})

}
