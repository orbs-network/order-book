package mocks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
)

// Mock store methods for service layer testing
type MockOrderBookStore struct {
	Error       error
	Order       *models.Order
	Orders      []models.Order
	MarketDepth models.MarketDepth
	OrderIter   models.OrderIter
	// auction
	Asks  []models.Order
	Bids  []models.Order
	Frags []models.OrderFrag
}

func (m *MockOrderBookStore) GetStore() service.OrderBookStore {
	return nil
}

func (m *MockOrderBookStore) StoreOrder(ctx context.Context, order models.Order) error {

	source, err := m.FindOrderById(ctx, order.Id, false)
	if err != nil {
		return err
	}

	//source.SizePending = order.SizePending
	*source = order

	return m.Error
}

func (m *MockOrderBookStore) StoreOrders(ctx context.Context, orders []models.Order) error {
	// update the orders
	for _, order := range orders {
		err := m.StoreOrder(ctx, order)
		if err != nil {
			return err
		}
	}
	return m.Error
}

func (m *MockOrderBookStore) RemoveOrder(ctx context.Context, order models.Order) error {
	return m.Error
}

func findOrder(orders []models.Order, id uuid.UUID) *models.Order {
	//for i, order := range *orders {
	for i := 0; i < len(orders); i++ { // := range *orders {
		if orders[i].Id == id {
			res := &orders[i]
			fmt.Printf("Address of Order %v", res)
			return res
		}
	}
	return nil
}
func (m *MockOrderBookStore) FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	order := findOrder(m.Asks, id)
	if order != nil {
		return order, nil
	}

	order = findOrder(m.Bids, id)
	if order != nil {
		return order, nil
	}

	order = findOrder(m.Orders, id)
	if order != nil {
		return order, nil
	}
	if m.Order == nil {
		return nil, models.ErrOrderNotFound
	}

	return m.Order, nil
}

func (m *MockOrderBookStore) GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Orders, nil
}

func (m *MockOrderBookStore) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (models.Order, error) {
	if m.Error != nil {
		return models.Order{}, m.Error
	}
	return *m.Order, nil
}

func (m *MockOrderBookStore) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	if m.Error != nil {
		return models.MarketDepth{}, m.Error
	}
	return m.MarketDepth, nil
}

func (m *MockOrderBookStore) StoreAuction(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error {
	if m.Error != nil {
		return m.Error
	}
	// save auction
	m.Frags = frags
	return nil
}

func (m *MockOrderBookStore) GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.Orders, len(m.Orders), nil
}

func (m *MockOrderBookStore) GetAuction(ctx context.Context, auctionID uuid.UUID) ([]models.OrderFrag, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Frags, nil
}

func (m *MockOrderBookStore) RemoveAuction(ctx context.Context, auctionID uuid.UUID) error {
	return nil
}

func (m *MockOrderBookStore) GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return m.OrderIter
}

func (m *MockOrderBookStore) GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return m.OrderIter
}

func (r *MockOrderBookStore) AddVal2Set(ctx context.Context, key, val string) error {
	return nil
}
