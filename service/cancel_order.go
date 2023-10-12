package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrder(ctx context.Context, orderId uuid.UUID) error {

	order, err := s.orderBookStore.FindOrderById(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "error occured when finding order", logger.Error(err))
		return err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("orderId", orderId.String()))
		return models.ErrOrderNotFound
	}

	err = s.orderBookStore.RemoveOrder(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "error occured when removing order", logger.Error(err))
		return err
	}

	logctx.Info(ctx, "order removed", logger.String("orderId", orderId.String()))
	return nil
}
