package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// Mock store methods for service layer testing
type MockOrderBookStore struct {
	Error error
	// Force a get/store user error
	ErrUser     error
	Order       *models.Order
	Orders      []models.Order
	User        *models.User
	MarketDepth models.MarketDepth
	OrderIter   models.OrderIter
	// auction
	Asks  []models.Order
	Bids  []models.Order
	Frags []models.OrderFrag
	// re-entrance
	Sets map[string]map[string]struct{}
}

func (m *MockOrderBookStore) FindOrdersByIds(ctx context.Context, ids []uuid.UUID) ([]models.Order, error) {
	return m.Orders, nil

}

func (m *MockOrderBookStore) StoreOpenOrder(ctx context.Context, order models.Order) error {

	source, err := m.FindOrderById(ctx, order.Id, false)
	if err != nil {
		return err
	}

	//source.SizePending = order.SizePending
	*source = order

	return m.Error
}

func (m *MockOrderBookStore) StoreOpenOrders(ctx context.Context, orders []models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) StoreFilledOrders(ctx context.Context, orders []models.Order) error {
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
		return nil, models.ErrNotFound
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

func (m *MockOrderBookStore) StoreSwap(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error {
	if m.Error != nil {
		return m.Error
	}
	// save auction
	m.Frags = frags
	return nil
}

func (m *MockOrderBookStore) GetOrdersForUser(ctx context.Context, userId uuid.UUID, isFilledOrders bool) (orders []models.Order, totalOrders int, err error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.Orders, len(m.Orders), nil
}

func (m *MockOrderBookStore) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	var orderIds []uuid.UUID
	for _, order := range m.Orders {
		orderIds = append(orderIds, order.Id)
	}
	return orderIds, m.Error
}

func (m *MockOrderBookStore) GetUserByPublicKey(ctx context.Context, publicKey string) (*models.User, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.User, nil
}

func (m *MockOrderBookStore) StoreUserByPublicKey(ctx context.Context, user models.User) error {
	return m.Error
}

func (m *MockOrderBookStore) GetAuction(ctx context.Context, auctionID uuid.UUID) ([]models.OrderFrag, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Frags, nil
}

func (m *MockOrderBookStore) RemoveAuction(ctx context.Context, auctionID uuid.UUID) error {
	m.Frags = []models.OrderFrag{}
	return nil
}

func (m *MockOrderBookStore) RemoveSwap(ctx context.Context, auctionID uuid.UUID) error {
	m.Frags = []models.OrderFrag{}
	return nil
}

func (m *MockOrderBookStore) GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return &OrderIterMock{
		orders: m.Asks,
		index:  -1,
	}
}

func (m *MockOrderBookStore) GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return &OrderIterMock{
		orders: m.Bids,
		index:  -1,
	}
}

func (m *MockOrderBookStore) UpdateAuctionTracker(ctx context.Context, auctionStatus models.AuctionStatus, auctionId uuid.UUID) error {
	return m.Error
}
