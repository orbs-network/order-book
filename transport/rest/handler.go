package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
)

// Service represents the methods available on the service to handle the actual request.
type Service interface {
	GetStore() service.OrderBookStore
	ProcessOrder(ctx context.Context, input service.ProcessOrderInput) (models.Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID, isClientOId bool) (cancelledOrderId *uuid.UUID, err error)
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error)
	GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (service.ConfirmAuctionRes, error)
	RevertAuction(ctx context.Context, auctionId uuid.UUID) error
	AuctionMined(ctx context.Context, auctionId uuid.UUID) error
	GetSymbols(ctx context.Context) ([]models.Symbol, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)
	CancelOrdersForUser(ctx context.Context, publicKey string) error
	GetAmountOut(ctx context.Context, auctionID uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error)
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
