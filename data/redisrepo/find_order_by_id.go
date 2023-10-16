package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) FindOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	orderIDKey := CreateOrderIDKey(orderId)

	orderMap, err := r.client.HGetAll(ctx, orderIDKey).Result()
	if err != nil {
		logctx.Error(ctx, "could not get order", logger.Error(err))
		return nil, err
	}

	if len(orderMap) == 0 {
		return nil, models.ErrOrderNotFound
	}

	order := &models.Order{}

	err = order.MapToOrder(orderMap)

	if err != nil {
		logctx.Error(ctx, "could not map order", logger.Error(err))
		return nil, err
	}

	return order, nil
}
