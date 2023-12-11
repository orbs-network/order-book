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
	StoreAuction(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error
	RemoveAuction(ctx context.Context, auctionID uuid.UUID) error
	GetAuction(ctx context.Context, auctionID uuid.UUID) ([]models.OrderFrag, error)
	GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter
	GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter
	UpdateAuctionTracker(ctx context.Context, auctionStatus models.AuctionStatus, auctionId uuid.UUID) error
	// taker side
	//UpdateAuctionTracker(ctx context.Context, auctionStatus models.AuctionStatus, auctionId uuid.UUID) error
	StoreSwap(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error
	RemoveSwap(ctx context.Context, auctionID uuid.UUID) error
}
