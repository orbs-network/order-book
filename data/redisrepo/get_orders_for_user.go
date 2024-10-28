package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// Fetches a user's open order IDs and their respective orders in one DB operation
func (r *redisRepository) GetOpenOrdersForUser(ctx context.Context, userId uuid.UUID) ([]models.Order, error) {
	pipeline := r.client.Pipeline()

	// Fetch all open order IDs for the user
	orderIdCmd := pipeline.ZRange(ctx, CreateUserOpenOrdersKey(userId), 0, -1)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute pipeline fetching order IDs", logger.Error(err))
		return nil, fmt.Errorf("failed to get order IDs for user: %v", err)
	}

	orderIdStrs, err := orderIdCmd.Result()
	if err != nil {
		logctx.Error(ctx, "failed to get order IDs result", logger.Error(err))
		return nil, fmt.Errorf("failed to get order IDs result: %v", err)
	}

	if len(orderIdStrs) == 0 {
		logctx.Warn(ctx, "no open orders found for user", logger.String("userId", userId.String()))
		return nil, models.ErrNotFound
	}

	// Get orders by IDs
	cmds := make([]*redis.MapStringStringCmd, len(orderIdStrs))
	for i, orderIdStr := range orderIdStrs {
		cmds[i] = pipeline.HGetAll(ctx, CreateOrderIDKey(uuid.MustParse(orderIdStr)))
	}

	_, err = pipeline.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute pipeline fetching orders by IDs", logger.Error(err), logger.Int("numIds", len(orderIdStrs)))
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	orders := make([]models.Order, 0, len(orderIdStrs))
	for i, cmd := range cmds {
		orderMap, err := cmd.Result()
		if err != nil {
			logctx.Error(ctx, "could not get order", logger.Error(err))
			return nil, fmt.Errorf("could not get order: %v", err)
		}

		if len(orderMap) == 0 {
			logctx.Warn(ctx, "order not found but was expected to exist", logger.String("orderId", orderIdStrs[i]))
			continue
		}

		order := models.Order{}
		err = order.MapToOrder(orderMap)
		if err != nil {
			logctx.Error(ctx, "could not map order", logger.Error(err))
			return nil, fmt.Errorf("could not map order: %v", err)
		}

		if order.IsOpen() {
			orders = append(orders, order)
		}
	}

	return orders, nil
}
