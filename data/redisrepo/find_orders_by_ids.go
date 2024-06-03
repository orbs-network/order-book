package redisrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

const MAX_ORDER_IDS = 500

// FindOrdersByIds finds orders by their IDs. If an order is not found for any of the provided IDs, an error is returned.
//
// Only finding orders by their IDs is supported, not by their clientOIds.

// Passing onlyOpen=true will only return orders that are open (not pending and not filled).
func (r *redisRepository) FindOrdersByIds(ctx context.Context, ids []uuid.UUID, onlyOpen bool) ([]models.Order, error) {

	if len(ids) > MAX_ORDER_IDS {
		return nil, fmt.Errorf("exceeded maximum number of IDs: %d", MAX_ORDER_IDS)
	}

	pipeline := r.client.Pipeline()

	cmds := make([]*redis.MapStringStringCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipeline.HGetAll(ctx, CreateOrderIDKey(id))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute pipeline getting orders by IDs", logger.Error(err), logger.Int("numIds", len(ids)))
		return nil, fmt.Errorf("failed to execute pipeline: %v", err)
	}

	orders := make([]models.Order, 0, len(ids))
	for _, cmd := range cmds {
		orderMap, err := cmd.Result()
		if err != nil {
			logctx.Error(ctx, "could not get order", logger.Error(err))
			return nil, fmt.Errorf("could not get order: %v", err)
		}

		if len(orderMap) == 0 {
			logctx.Warn(ctx, "order not found but was expected to exist", logger.String("orderId", cmd.Args()[1].(string)))
			return nil, errors.New("order not found but was expected to exist")
		}

		order := models.Order{}
		err = order.MapToOrder(orderMap)
		if err != nil {
			logctx.Error(ctx, "could not map order", logger.Error(err))
			return nil, fmt.Errorf("could not map order: %v", err)
		}

		if onlyOpen {
			if order.IsOpen() {
				orders = append(orders, order)
			}
		} else {
			orders = append(orders, order)
		}
	}

	logctx.Debug(ctx, "found orders by IDs", logger.Int("numIdsProvided", len(ids)), logger.Int("numOrders", len(orders)))
	return orders, nil
}
