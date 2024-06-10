package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error) {
	order, err := s.orderBookStore.FindOrderById(ctx, clientOId, true)

	if err == models.ErrNotFound {
		logctx.Debug(ctx, "order not found", logger.String("clientOId", clientOId.String()))
		return nil, nil
	}

	if err != nil {
		logctx.Error(ctx, "unexpected error when getting order", logger.String("clientOId", clientOId.String()), logger.Error(err))
		return nil, err
	}

	return order, nil
}
