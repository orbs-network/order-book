package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetOpenOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	logctx.Info(ctx, "getting open orders for user", logger.String("user_id", userId.String()))
	orders, totalOrders, err = s.orderBookStore.GetOrdersForUser(ctx, userId, false)

	if err != nil {
		logctx.Error(ctx, "error getting open orders for user", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, 0, fmt.Errorf("error getting open orders for user: %w", err)
	}

	logctx.Info(ctx, "returning open orders for user", logger.String("user_id", userId.String()), logger.Int("orders_count", len(orders)))

	return orders, totalOrders, nil
}

func (s *Service) GetFilledOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	logctx.Info(ctx, "getting filled orders for user", logger.String("user_id", userId.String()))
	orders, totalOrders, err = s.orderBookStore.GetOrdersForUser(ctx, userId, true)

	if err != nil {
		logctx.Error(ctx, "error getting filled orders for user", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, 0, fmt.Errorf("error getting filled orders for user: %w", err)
	}

	logctx.Info(ctx, "returning filled orders for user", logger.String("user_id", userId.String()), logger.Int("orders_count", len(orders)))

	return orders, totalOrders, nil
}
