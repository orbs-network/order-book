package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	order, err := s.orderBookStore.FindOrderById(ctx, orderId, false)

	if err == models.ErrNotFound {
		logctx.Info(ctx, "order not found", logger.String("orderId", orderId.String()))
		return nil, nil
	}

	if err != nil {
		logctx.Error(ctx, "unexpected error when getting order", logger.String("orderId", orderId.String()), logger.Error(err))
		return nil, err
	}

	logctx.Info(ctx, "order found for orderId", logger.String("orderId", orderId.String()))
	return order, nil
}
