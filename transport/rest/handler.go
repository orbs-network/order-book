package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// Service represents the methods available on the service to handle the actual request.
type Service interface {
	// TODO: rename to ProcessOrder as sometimes an order will be immediately filled and not added at all to order book
	AddOrder(ctx context.Context, userId uuid.UUID, price decimal.Decimal, symbol models.Symbol, size decimal.Decimal, side models.Side) (models.Order, error)
	CancelOrder(ctx context.Context, orderId uuid.UUID) error
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
}

type Handler struct {
	svc    Service
	router *mux.Router
}

func NewHandler(svc Service, r *mux.Router) (*Handler, error) {
	if svc == nil {
		return nil, fmt.Errorf("svc cannot be nil")
	}

	if r == nil {
		return nil, fmt.Errorf("router cannot be nil")
	}

	return &Handler{
		svc:    svc,
		router: r,
	}, nil
}
