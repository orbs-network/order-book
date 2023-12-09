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

// GetOrdersForUser returns all orders (open or filled) for a given user, sorted by creation time.
//
// `isFilledOrders` should be true if you want filled orders, false if you want open orders.
//
// This function is paginated, and returns the total number of orders for the user
func (r *redisRepository) GetOrdersForUser(ctx context.Context, userId uuid.UUID, isFilledOrders bool) (orders []models.Order, totalOrders int, err error) {
	start, stop := utils.PaginationBounds(ctx)

	var key string
	if isFilledOrders {
		logctx.Info(ctx, "getting filled orders for user", logger.String("userId", userId.String()))
		key = CreateUserFilledOrdersKey(userId)
	} else {
		logctx.Info(ctx, "getting open orders for user", logger.String("userId", userId.String()))
		key = CreateUserOpenOrdersKey(userId)
	}

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

	orderIds := make([]uuid.UUID, len(orderIdStrs))
	for i, orderIdStr := range orderIdStrs {
		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "failed to parse order id", logger.String("orderId", orderIdStr), logger.Error(err))
			return []models.Order{}, 0, fmt.Errorf("failed to parse order id: %w", err)
		}
		orderIds[i] = orderId
	}

	orders, err = r.FindOrdersByIds(ctx, orderIds)
	if err != nil {
		logctx.Error(ctx, "failed to get orders", logger.String("userId", userId.String()), logger.Error(err))
		return []models.Order{}, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	logctx.Info(ctx, "got orders for user", logger.String("userId", userId.String()), logger.Int("count", len(orders)))
	return orders, totalOrders, nil
}
