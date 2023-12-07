package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) (orderIds []uuid.UUID, err error) {

	orderIds, err = s.orderBookStore.CancelOrdersForUser(ctx, userId)

	if err == models.ErrNoOrdersFound {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return []uuid.UUID{}, err
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", userId.String()))
		return []uuid.UUID{}, fmt.Errorf("could not cancel orders for user: %w", err)
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", userId.String()), logger.Int("numOrders", len(orderIds)))
	return orderIds, nil
}
