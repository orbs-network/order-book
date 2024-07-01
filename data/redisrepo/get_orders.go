package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// returns all open IDs if symbol is empty
func (r *redisRepository) GetOpenOrderIds(ctx context.Context, userId uuid.UUID, symbol models.Symbol) ([]uuid.UUID, error) {
	userOrdersKey := CreateUserOpenOrdersKey(userId)

	// Fetch all order IDs for the user
	orderIdStrs, err := r.client.ZRange(ctx, userOrdersKey, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get order IDs for user", logger.String("userId", userId.String()), logger.Error(err))
		return nil, fmt.Errorf("failed to get order IDs for user: %v", err)
	}

	if len(orderIdStrs) == 0 {
		logctx.Warn(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return nil, models.ErrNotFound
	}
	// Convert string IDs to UUIDs
	var orderIds []uuid.UUID
	for _, orderIdStr := range orderIdStrs {
		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "failed to parse order ID", logger.String("orderId", orderIdStr), logger.Error(err))
			return nil, fmt.Errorf("failed to parse order ID: %v", err)
		}
		orderIds = append(orderIds, orderId)
	}
	return orderIds, nil

}
