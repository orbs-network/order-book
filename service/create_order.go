package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CreateOrder(ctx context.Context, price string, symbol string, size string) (models.Order, error) {
	id := uuid.New().String()
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

	logctx.Info(ctx, "order added", logger.String("ID", order.Id), logger.String("price", order.Price), logger.String("size", order.Size))
	return order, nil
}
