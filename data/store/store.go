package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type OrderBookStore interface {
	// MM side
	StoreOpenOrder(ctx context.Context, order models.Order) error
	StoreOpenOrders(ctx context.Context, orders []models.Order) error
	StoreFilledOrders(ctx context.Context, orders []models.Order) error
	RemoveOrder(ctx context.Context, order models.Order) error
	FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error)
	FindOrdersByIds(ctx context.Context, ids []uuid.UUID) ([]models.Order, error)
	GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error)
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID, isFilledOrders bool) (orders []models.Order, totalOrders int, err error)
	CancelOrdersForUser(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	// LH side
	GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter
	GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter
	// taker side
	UpdateSwapTracker(ctx context.Context, swapStatus models.SwapStatus, swapId uuid.UUID) error
	GetSwap(ctx context.Context, swapId uuid.UUID) ([]models.OrderFrag, error)
	StoreSwap(ctx context.Context, swapId uuid.UUID, frags []models.OrderFrag) error
	RemoveSwap(ctx context.Context, swapId uuid.UUID) error
	// Pending transactions
	StoreNewPendingSwap(ctx context.Context, pendingSwap models.Pending) error
	GetPendingSwaps(ctx context.Context) ([]models.Pending, error)
	StorePendingSwaps(ctx context.Context, pendingSwaps []models.Pending) error
	ProcessCompletedSwapOrders(ctx context.Context, orders []*models.Order, swapId uuid.UUID, isSuccessful bool) error
}
