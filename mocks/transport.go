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
	AmountOut   models.AmountOut
}

func (m *MockOrderBookService) ProcessOrder(ctx context.Context, input service.ProcessOrderInput) (models.Order, error) {
	return *m.Order, m.Error
}

func (m *MockOrderBookService) CancelOrder(ctx context.Context, id uuid.UUID, isClientOId bool) (cancelledOrderId *uuid.UUID, err error) {
	if m.Error != nil {
		return nil, m.Error
	}

	if m.Order == nil {
		return nil, nil
	}

	return &m.Order.Id, m.Error
}

func (m *MockOrderBookService) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error) {
	return decimal.Zero, m.Error
}

func (m *MockOrderBookService) GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	return m.Order, m.Error
}

func (m *MockOrderBookService) GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error) {
	return m.Order, m.Error
}

func (m *MockOrderBookService) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	return m.MarketDepth, m.Error
}

func (m *MockOrderBookService) GetAmountOut(ctx context.Context, auctionId string, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {
	return m.AmountOut, m.Error
}
func (m *MockOrderBookService) ConfirmAuction(ctx context.Context, auctionId string) (service.ConfirmAuctionRes, error) {
	return service.ConfirmAuctionRes{}, nil
}
func (m *MockOrderBookService) RemoveAuction(ctx context.Context, auctionId string) error {
	return nil
}
