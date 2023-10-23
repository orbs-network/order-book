package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
)

// Mock service methods for transport layer testing
type MockOrderBookService struct {
	Error       error
	Order       *models.Order
	Orders      []models.Order
	MarketDepth models.MarketDepth
}

func (m *MockOrderBookService) ProcessOrder(ctx context.Context, input service.ProcessOrderInput) (models.Order, error) {
	return *m.Order, m.Error
}

func (m *MockOrderBookService) CancelOrder(ctx context.Context, orderId uuid.UUID) error {
	return m.Error
}

func (m *MockOrderBookService) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error) {
	return decimal.Zero, m.Error
}

func (m *MockOrderBookService) GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	return m.Order, m.Error
}

func (m *MockOrderBookService) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	return m.MarketDepth, m.Error
}
