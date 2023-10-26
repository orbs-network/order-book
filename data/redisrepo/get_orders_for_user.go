package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// TODO: add support for bulk getting orders by IDs

// GetOrdersForUser returns all orders for a given user, sorted by creation time.
// This function is paginated, and returns the total number of orders for the user
func (r *redisRepository) GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error) {
	start, stop := utils.PaginationBounds(ctx)

	key := CreateUserOrdersKey(userId)

	count, err := r.client.ZCard(ctx, key).Result()

	if err != nil {
		logctx.Error(ctx, "failed to get total count of orders for user", logger.String("userId", userId.String()), logger.Error(err))
		return []models.Order{}, 0, fmt.Errorf("failed to get total count of orders for user: %w", err)
	}

	totalOrders = int(count)

	orderIdStrs, err := r.client.ZRange(ctx, key, start, stop).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get orders for user", logger.String("userId", userId.String()), logger.Error(err))
		return []models.Order{}, 0, fmt.Errorf("failed to get orders for user: %w", err)
	}

	orders = make([]models.Order, len(orderIdStrs))

	for i, orderIdStr := range orderIdStrs {
		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "failed to parse order id", logger.String("orderId", orderIdStr), logger.Error(err))
			return []models.Order{}, 0, fmt.Errorf("failed to parse order id: %w", err)
		}

		order, err := r.FindOrderById(ctx, orderId, false)
		if err != nil {
			logctx.Error(ctx, "failed to get order", logger.String("orderId", orderIdStr), logger.Error(err))
			return []models.Order{}, 0, fmt.Errorf("failed to get order: %w", err)
		}
		if order == nil {
			logctx.Error(ctx, "order not found but should exist", logger.String("orderId", orderIdStr))
			return []models.Order{}, 0, fmt.Errorf("order not found but should exist")
		}

		orders[i] = *order
	}

	logctx.Info(ctx, "got orders for user", logger.String("userId", userId.String()), logger.Int("count", len(orders)))
	return orders, totalOrders, nil
}
