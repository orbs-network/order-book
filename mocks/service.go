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
	ErrUser      error
	Order        *models.Order
	Orders       []models.Order
	User         *models.User
	MarketDepth  models.MarketDepth
	AskOrderIter models.OrderIter
	BidOrderIter models.OrderIter
	// swap
	Asks  []models.Order
	Bids  []models.Order
	Frags []models.OrderFrag
	// re-entrance
	Sets map[string]map[string]struct{}
	// Pending swaps
	PendingSwaps []models.SwapTx
	// PubSub
	EventsChan chan []byte
}

func (m *MockOrderBookStore) StoreOpenOrder(ctx context.Context, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) StoreOpenOrders(ctx context.Context, orders []models.Order) error {

	return m.Error
}

func (m *MockOrderBookStore) StoreFilledOrders(ctx context.Context, orders []models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) CancelPartialFilledOrder(ctx context.Context, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) FindOrdersByIds(ctx context.Context, ids []uuid.UUID, onlyOpen bool) ([]models.Order, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Orders, nil
}

func (m *MockOrderBookStore) FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error) {

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

func (m *MockOrderBookStore) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	if m.Error != nil {
		return models.MarketDepth{}, m.Error
	}
	return m.MarketDepth, nil
}

func (m *MockOrderBookStore) StoreSwap(ctx context.Context, swapId uuid.UUID, symbol models.Symbol, side models.Side, frags []models.OrderFrag) error {
	if m.Error != nil {
		return m.Error
	}

	return nil
}

func (m *MockOrderBookStore) GetOpenOrders(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orders []models.Order, totalOrders int, err error) {
	if m.Error != nil {
		return nil, 0, m.Error
	}
	return m.Orders, len(m.Orders), nil
}

func (m *MockOrderBookStore) GetOpenOrdersForUser(ctx context.Context, userId uuid.UUID) ([]models.Order, error) {
	return m.Orders, m.Error
}

// Generic Building blocks with no biz logic in a single TX
func (m *MockOrderBookStore) PerformTx(ctx context.Context, action func(txid uint) error) error {
	return m.Error
}

func (m *MockOrderBookStore) TxModifyOrder(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) TxModifyPrices(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) TxModifyClientOId(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) TxModifyUserOpenOrders(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) TxRemoveOrder(ctx context.Context, txid uint, order models.Order) error {
	return m.Error
}

func (m *MockOrderBookStore) TxCloseOrder(ctx context.Context, txid uint, order models.Order) error {
	return m.Error
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

func (m *MockOrderBookStore) GetSwap(ctx context.Context, swapId uuid.UUID, open bool) (*models.Swap, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return models.NewSwap(models.Symbol("MATIC_USDC"), models.BUY, m.Frags), nil
}

func (m *MockOrderBookStore) RemoveSwap(ctx context.Context, swapId uuid.UUID) error {

	return m.Error
}

func (m *MockOrderBookStore) GetOpenSwaps(ctx context.Context) ([]models.Swap, error) {
	return []models.Swap{}, m.Error
}

func (m *MockOrderBookStore) GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return m.AskOrderIter
}

func (m *MockOrderBookStore) GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter {
	return m.BidOrderIter
}

func (m *MockOrderBookStore) StoreNewPendingSwap(ctx context.Context, pendingSwap models.SwapTx) (*models.Swap, error) {
	swap := models.NewSwap(m.Order.Symbol, m.Order.Side, m.Frags)
	return swap, m.Error
}

func (m *MockOrderBookStore) ResolveSwap(ctx context.Context, swap models.Swap) error {
	return m.Error
}

func (m *MockOrderBookStore) StoreUserResolvedSwap(ctx context.Context, userId uuid.UUID, swap models.Swap) error {
	return m.Error
}

func (m *MockOrderBookStore) GetUserResolvedSwapIds(ctx context.Context, userId uuid.UUID) ([]string, error) {
	return []string{"111", "222"}, m.Error
}

func (m *MockOrderBookStore) EnumSubKeysOf(tx context.Context, key string) ([]string, error) {
	return []string{key + "111", key + "222"}, m.Error
}

func (m *MockOrderBookStore) PublishEvent(ctx context.Context, key string, value interface{}) error {
	return m.Error
}

func (m *MockOrderBookStore) SubscribeToEvents(ctx context.Context, channel string) (chan []byte, error) {
	return m.EventsChan, m.Error
}

func (m *MockOrderBookStore) UnsubscribeFromEvents(ctx context.Context, channel string, clientChan chan []byte) {
}

func (r *MockOrderBookStore) ReadStrKey(ctx context.Context, key string) (string, error) {
	return "", nil
}
func (r *MockOrderBookStore) WriteStrKey(ctx context.Context, key, val string) error {
	return nil
}

func (r *MockOrderBookStore) GetMakerTokenBalance(ctx context.Context, token, wallet string) (decimal.Decimal, error) {
	return decimal.Zero, nil
}
