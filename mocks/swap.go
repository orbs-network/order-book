package mocks

import (
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func newOrder(price, size int64) models.Order {
	oid, _ := uuid.NewUUID()
	return models.Order{
		Id:          oid,
		Price:       decimal.NewFromInt(price),
		Size:        decimal.NewFromInt(size),
		SizePending: decimal.Zero,
		SizeFilled:  decimal.Zero,
	}
}

func newAsks() []models.Order {
	return []models.Order{
		newOrder(1000, 1),
		newOrder(1001, 2),
		newOrder(1002, 3),
	}
}
func newBids() []models.Order {
	return []models.Order{
		newOrder(900, 1),
		newOrder(800, 2),
		newOrder(700, 3),
	}
}

func newFrags(orders []models.Order) []models.OrderFrag {
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

func CreateSwapMock() *MockOrderBookStore {

	asks := newAsks()
	bids := newBids()

	askOrderIter := OrderIterMock{
		Orders: asks,
		Index:  -1,
	}

	bidOrderIter := OrderIterMock{
		Orders: bids,
		Index:  -1,
	}

	res := MockOrderBookStore{
		Error:        nil,
		Sets:         make(map[string]map[string]struct{}),
		AskOrderIter: &askOrderIter,
		BidOrderIter: &bidOrderIter,
		Order:        &bids[0],
	}
	res.Asks = asks
	res.Bids = bids
	res.Frags = newFrags(res.Asks)
	return &res
}
