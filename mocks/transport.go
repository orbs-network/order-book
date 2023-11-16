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
	Symbols     []models.Symbol
	User        *models.User
}

func (m *MockOrderBookService) GetUserByPublicKey(ctx context.Context, publicKey string) (*models.User, error) {
	return m.User, m.Error
}

func (m *MockOrderBookService) ProcessOrder(ctx context.Context, input service.ProcessOrderInput) (models.Order, error) {
	return *m.Order, m.Error
}

func (m *MockOrderBookService) CancelOrder(ctx context.Context, input service.CancelOrderInput) (cancelledOrderId *uuid.UUID, err error) {
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

func (m *MockOrderBookService) GetSymbols(ctx context.Context) ([]models.Symbol, error) {
	return m.Symbols, m.Error
}

func (m *MockOrderBookService) GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	return m.Orders, len(m.Orders), m.Error
}

func (m *MockOrderBookService) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) error {
	return m.Error
}

func (m *MockOrderBookService) GetAmountOut(ctx context.Context, auctionId uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {
	return m.AmountOut, m.Error
}
func (m *MockOrderBookService) ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (service.ConfirmAuctionRes, error) {
	return service.ConfirmAuctionRes{}, nil
}
func (m *MockOrderBookService) RevertAuction(ctx context.Context, auctionId uuid.UUID) error {
	return nil
}

func (m *MockOrderBookService) AuctionMined(ctx context.Context, auctionId uuid.UUID) error {
	return nil
}
