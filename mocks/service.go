package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// Mock store methods for service layer testing
type MockOrderBookStore struct {
	Error       error
	Order       *models.Order
	Orders      []models.Order
	MarketDepth models.MarketDepth
}

func (m *MockOrderBookStore) StoreOrder(ctx context.Context, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) RemoveOrder(ctx context.Context, orderId uuid.UUID) error {
	return m.Error
}

func (m *MockOrderBookStore) FindOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	if m.Error != nil {
		return nil, m.Error
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
