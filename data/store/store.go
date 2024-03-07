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
	// MM side
	StoreOpenOrder(ctx context.Context, order models.Order) error
	// store multiple orders in a single redis tx in every state unfilled/pending/partial
	StoreOpenOrders(ctx context.Context, orders []models.Order) error
	// store multiple FILLED orders in a single redis tx
	// adds to user's filled orders USERID:filledOrders
	// removes orders from price list, removes from user's open orders and
	StoreFilledOrders(ctx context.Context, orders []models.Order) error
	// Order is completely removed from DB
	// Order is removed from the prices sorted set, user's open order set and order hash is removed
	// May only be called if order is not pending and completely unfilled
	CancelUnfilledOrder(ctx context.Context, order models.Order) error
	// Order remains in DB, but is marked as cancelled
	// Order is removed from the prices sorted set, user's order set and order hash is updated to cancelled
	// May only be called if order is not pending and partially filled
	CancelPartialFilledOrder(ctx context.Context, order models.Order) error
	// Order remains in DB, but is marked as cancelled
	// Order is removed from the prices sorted set, user's order set and order hash is updated to cancelled
	// Upon swap resolve false -> should be removed
	CancelPendingOrder(ctx context.Context, order models.Order) error
	FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error)
	FindOrdersByIds(ctx context.Context, ids []uuid.UUID, onlyOpen bool) ([]models.Order, error)
	GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID, isFilledOrders bool) (orders []models.Order, totalOrders int, err error)
	CancelOrdersForUser(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	// Generic Building blocks with no biz logic in a single TX
	TxStart(ctx context.Context) (uint, error)
	TxEnd(ctx context.Context, txid uint)
	TxRemoveOrderFromPrice(ctx context.Context, txid uint, order models.Order) error
	TxDeleteOrder(ctx context.Context, txid uint, orderId uuid.UUID) error
	TxStoreOrder(ctx context.Context, txid uint, order models.Order) error

	// LH side
	GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter
	GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter
	// taker side
	GetSwap(ctx context.Context, swapId uuid.UUID) (*models.Swap, error)
	StoreSwap(ctx context.Context, swapId uuid.UUID, frags []models.OrderFrag) error
	RemoveSwap(ctx context.Context, swapId uuid.UUID) error
	GetOpenSwaps(ctx context.Context) ([]models.Swap, error)
	// Pending Swap+Transaction (TODO: rename)
	StoreNewPendingSwap(ctx context.Context, pendingSwap models.SwapTx) error
	// removes from "swapid" key
	// adds to "swapResolve" key
	ResolveSwap(ctx context.Context, swap models.Swap) error
	// save swapId in a set of the userId:resolvedSwap key
	StoreUserResolvedSwap(ctx context.Context, userId uuid.UUID, swap models.Swap) error
	// utils
	EnumSubKeysOf(ctx context.Context, key string) ([]string, error)
}
