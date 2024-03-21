// Package service contains the business logic for the application.

package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type OrderBookService interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (models.Order, error)
	CancelOrder(ctx context.Context, input CancelOrderInput) (*uuid.UUID, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error)
	GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	CancelOrdersForUser(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orderIds []uuid.UUID, err error)
	GetSymbols(ctx context.Context) ([]models.Symbol, error)
	GetOpenOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)
	GetFilledOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)

	// taker api - INSTEAD
	GetQuote(ctx context.Context, symbol models.Symbol, side models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal) (models.QuoteRes, error)
	BeginSwap(ctx context.Context, data models.QuoteRes) (models.BeginSwapRes, error)
	SwapStarted(ctx context.Context, swapId uuid.UUID, txHash string) error
	AbortSwap(ctx context.Context, swapId uuid.UUID) error
	FillSwap(ctx context.Context, swapId uuid.UUID) error
}

type BlockChainService interface {
	CheckPendingTxs(ctx context.Context) error
}

// Service contains methods that implement the business logic for the application.
type Service struct {
	orderBookStore   store.OrderBookStore
	blockchainClient BlockChainService
	reporter         *Reporter
}

// New creates a new Service with injected dependencies.
func New(store store.OrderBookStore, bcClient BlockChainService) (*Service, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	if bcClient == nil {
		return nil, errors.New("bcClient cannot be nil")
	}

	// start report routine
	svc := Service{orderBookStore: store, blockchainClient: bcClient}
	svc.reporter = NewReporter(&svc)
	svc.reporter.Start()

	// start periodic check routine
	svc.startPeriodicChecks()

	return &svc, nil
}
