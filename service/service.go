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
	ProcessOrder(ctx context.Context, input ProcessOrderInput) (models.Order, error)
	CancelOrder(ctx context.Context, input CancelOrderInput) (cancelledOrderId *uuid.UUID, err error)
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error)
	GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	CancelOrdersForUser(ctx context.Context, userId uuid.UUID) (orderIds []uuid.UUID, err error)
	GetSymbols(ctx context.Context) ([]models.Symbol, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)

	ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (ConfirmAuctionRes, error)
	RevertAuction(ctx context.Context, auctionId uuid.UUID) error
	AuctionMined(ctx context.Context, auctionId uuid.UUID) error
	GetAmountOut(ctx context.Context, auctionID uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error)
}

type BlockChainService interface {
	VerifySignature(ctx context.Context, input VerifySignatureInput) (bool, error)
}

// Service contains methods that implement the business logic for the application.
type Service struct {
	orderBookStore   store.OrderBookStore
	blockchainClient BlockChainService
}

// New creates a new Service with injected dependencies.
func New(store store.OrderBookStore, bcClient BlockChainService) (*Service, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	if bcClient == nil {
		return nil, errors.New("bcClient cannot be nil")
	}

	return &Service{orderBookStore: store, blockchainClient: bcClient}, nil
}
