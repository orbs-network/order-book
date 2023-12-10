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

// FindOrderById finds an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOId, otherwise it is treated as an orderId.
func (r *redisRepository) FindOrderById(ctx context.Context, id uuid.UUID, isClientOId bool) (*models.Order, error) {
	var orderId uuid.UUID

	if isClientOId {
		logctx.Info(ctx, "finding order by clientOId", logger.String("clientOId", id.String()))

		orderIdStr, err := r.client.Get(ctx, CreateClientOIDKey(id)).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, models.ErrNotFound
			}
			logctx.Error(ctx, "could not get order ID by clientOId", logger.Error(err))
			return nil, err
		}

		orderId, err = uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "invalid order ID format retrieved by clientOId", logger.Error(err))
			return nil, fmt.Errorf("invalid order ID format retrieved by clientOId: %s", err)
		}
	} else {
		logctx.Info(ctx, "finding order by orderId", logger.String("orderId", id.String()))
		orderId = id
	}

	orderMap, err := r.client.HGetAll(ctx, CreateOrderIDKey(orderId)).Result()
	if err != nil {
		logctx.Error(ctx, "could not get order", logger.Error(err))
		return nil, err
	}

	if len(orderMap) == 0 {
		return nil, models.ErrNotFound
	}

	order := &models.Order{}
	err = order.MapToOrder(orderMap)
	if err != nil {
		logctx.Error(ctx, "could not map order", logger.Error(err))
		return nil, err
	}

	logctx.Info(ctx, "found order", logger.String("orderId", order.Id.String()))
	return order, nil
}
