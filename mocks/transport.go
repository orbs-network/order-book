package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
)

// Mock service methods for transport layer testing
type MockOrderBookService struct {
	Error        error
	Order        *models.Order
	Orders       []models.Order
	MarketDepth  models.MarketDepth
	QuoteRes     models.QuoteRes
	Symbols      []models.Symbol
	User         *models.User
	BeginSwapRes models.BeginSwapRes
	OrderEvents  chan []byte
}

func (m *MockOrderBookService) GetUserByPublicKey(ctx context.Context, publicKey string) (*models.User, error) {
	return m.User, m.Error
}

func (m *MockOrderBookService) CreateOrder(ctx context.Context, input service.CreateOrderInput) (models.Order, error) {
	return *m.Order, m.Error
}

func (m *MockOrderBookService) CancelOrder(ctx context.Context, input service.CancelOrderInput) (*uuid.UUID, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	if m.Order == nil {
		return nil, nil
	}

	return &m.Order.Id, m.Error
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

func (m *MockOrderBookService) GetOpenOrders(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orders []models.Order, totalOrders int, err error) {
	return m.Orders, len(m.Orders), m.Error
}

func (m *MockOrderBookService) GetSwapFills(ctx context.Context, userId uuid.UUID, symbol models.Symbol, startAt, endAt time.Time) ([]models.Fill, error) {
	return nil, m.Error
}

func (m *MockOrderBookService) CancelOrdersForUser(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orderIds []uuid.UUID, err error) {
	var ids []uuid.UUID
	for _, order := range m.Orders {
		ids = append(ids, order.Id)
	}
	return ids, m.Error
}

func (m *MockOrderBookService) SubscribeUserOrders(ctx context.Context, userId uuid.UUID) (chan []byte, error) {
	return m.OrderEvents, m.Error
}

func (m *MockOrderBookService) UnsubscribeUserOrders(ctx context.Context, userId uuid.UUID, clientChan chan []byte) error {
	return m.Error
}

func (m *MockOrderBookService) GetQuote(ctx context.Context, symbol models.Symbol, makerSide models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal, makerInToken string) (models.QuoteRes, error) {
	return m.QuoteRes, m.Error
}

// taker api instead of swap
func (m *MockOrderBookService) BeginSwap(ctx context.Context, data models.QuoteRes) (models.BeginSwapRes, error) {
	return m.BeginSwapRes, m.Error
}

func (m *MockOrderBookService) AbortSwap(ctx context.Context, swapId uuid.UUID) error {
	return m.Error
}

func (m *MockOrderBookService) SwapStarted(ctx context.Context, swapId uuid.UUID, txHash string) error {
	return m.Error
}

func (m *MockOrderBookService) FillSwap(ctx context.Context, swapId uuid.UUID) error {
	return m.Error
}
