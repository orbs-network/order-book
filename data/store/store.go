package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type OrderWithSize struct {
	Order *models.Order
	Size  decimal.Decimal
}

type OrderBookStore interface {
	StoreOpenOrder(ctx context.Context, order models.Order) error
	StoreOpenOrders(ctx context.Context, orders []models.Order) error
	// --- MM side ---
	// DEPRECATED - Use TxModify__ methods instead
	StoreFilledOrders(ctx context.Context, orders []models.Order) error
	FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error)
	FindOrdersByIds(ctx context.Context, ids []uuid.UUID, onlyOpen bool) ([]models.Order, error)
	GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	GetOpenOrders(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orders []models.Order, totalOrders int, err error)
	// ------------------------------
	// Generic getters
	GetOpenOrderIds(ctx context.Context, userId uuid.UUID, symbol models.Symbol) ([]uuid.UUID, error)
	// ------------------------------
	// Generic Building blocks with no biz logic in a single tx

	// PerformTX should be used for all interactions with the Redis repository. Handles the transaction lifecycle.
	PerformTx(ctx context.Context, action func(txid uint) error) error
	TxModifyOrder(ctx context.Context, txid uint, operation models.Operation, order models.Order) error
	TxModifyPrices(ctx context.Context, txid uint, operation models.Operation, order models.Order) error
	TxModifyClientOId(ctx context.Context, txid uint, operation models.Operation, order models.Order) error
	TxModifyUserOpenOrders(ctx context.Context, txid uint, operation models.Operation, order models.Order) error
	TxCloseOrder(ctx context.Context, txid uint, order models.Order) error
	TxRemoveOrder(ctx context.Context, txid uint, order models.Order) error
	// ------------------------------
	// LH side
	GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter
	GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter
	// taker side
	GetSwap(ctx context.Context, swapId uuid.UUID, open bool) (*models.Swap, error)
	StoreSwap(ctx context.Context, swapId uuid.UUID, symbol models.Symbol, side models.Side, frags []models.OrderFrag) error
	RemoveSwap(ctx context.Context, swapId uuid.UUID) error
	GetOpenSwaps(ctx context.Context) ([]models.Swap, error)
	// Pending Swap+Transaction (TODO: rename)
	StoreNewPendingSwap(ctx context.Context, pendingSwap models.SwapTx) (*models.Swap, error)
	// removes from "swapid" key
	// adds to "swapResolve" key
	ResolveSwap(ctx context.Context, swap models.Swap) error
	// save swapId in a set of the userId:resolvedSwap key
	StoreUserResolvedSwap(ctx context.Context, userId uuid.UUID, swap models.Swap) error
	GetUserResolvedSwapIds(ctx context.Context, userId uuid.UUID) ([]string, error)

	// utils
	EnumSubKeysOf(ctx context.Context, key string) ([]string, error)
	ReadStrKey(ctx context.Context, key string) (string, error)
	WriteStrKey(ctx context.Context, key, val string) error
	GetMakerTokenBalance(ctx context.Context, token, wallet string) (decimal.Decimal, error)

	// PubSub
	PublishEvent(ctx context.Context, key string, value interface{}) error
	SubscribeToEvents(ctx context.Context, channel string) (chan []byte, error)
}
