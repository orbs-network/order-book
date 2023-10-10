package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateOrder(ctx context.Context, price decimal.Decimal, symbol models.Symbol, size decimal.Decimal) (models.Order, error) {
	id := uuid.New()
	order := models.Order{
		Id:        id,
		Price:     price,
		Symbol:    symbol,
		Size:      size,
		Signature: nil,
		Pending:   false,
	}

	order, err := s.store.AddOrder(ctx, order)
	if err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "order added", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))
	return order, nil
}
