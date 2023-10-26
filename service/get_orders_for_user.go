package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	orders, totalOrders, err = s.orderBookStore.GetOrdersForUser(ctx, userId)

	if err != nil {
		logctx.Error(ctx, "error getting orders for user", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, 0, fmt.Errorf("error getting orders for user: %w", err)
	}

	logctx.Info(ctx, "returning orders for user", logger.String("user_id", userId.String()), logger.Int("orders_count", len(orders)))

	return orders, totalOrders, nil
}
