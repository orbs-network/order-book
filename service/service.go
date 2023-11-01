// Package service contains the business logic for the application.

package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type OrderBookStore interface {
	// MM side
	StoreOrder(ctx context.Context, order models.Order) error
	StoreOrders(ctx context.Context, orders []models.Order) error
	RemoveOrder(ctx context.Context, order models.Order) error
	FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error)
	GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error)
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)
	// LH side
	StoreAuction(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error
	RemoveAuction(ctx context.Context, auctionID uuid.UUID) error
	GetAuction(ctx context.Context, auctionID uuid.UUID) ([]models.OrderFrag, error)
	GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter
	GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter
}

// Service contains methods that implement the business logic for the application.
type Service struct {
	orderBookStore OrderBookStore
}

// New creates a new Service with injected dependencies.
func New(store OrderBookStore) (*Service, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	return &Service{orderBookStore: store}, nil
}
