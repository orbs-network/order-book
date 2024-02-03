package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// TODO: Refactor this.
//
// TODO 1 - Introduce transaction support in service layer.
//
// TODO 2 -. Do this logic in the service layer. Too much business logic here.
func (r *redisRepository) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
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

	// Fetch all orders for the user (in batches)
	var ordersToCancel []models.Order
	for i := 0; i < len(orderIds); i += MAX_ORDER_IDS {
		end := i + MAX_ORDER_IDS
		if end > len(orderIds) {
			end = len(orderIds)
		}

		// We only want to fetch open orders
		orders, err := r.FindOrdersByIds(ctx, orderIds[i:end], true)
		if err != nil {
			logctx.Error(ctx, "failed to find orders by IDs", logger.String("userId", userId.String()), logger.Error(err))
			return nil, fmt.Errorf("failed to find orders by IDs: %v", err)
		}

		ordersToCancel = append(ordersToCancel, orders...)
	}

	for _, order := range ordersToCancel {

		if order.IsUnfilled() {
			r.CancelUnfilledOrder(ctx, order)
			continue
		}

		if order.IsPartialFilled() {
			r.CancelPartialFilledOrder(ctx, order)
			continue
		}

		logctx.Error(ctx, "encountered order in an unexpected state when cancelling all orders for user", logger.String("orderId", order.Id.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))
	}

	logctx.Info(ctx, "removed all orders for user", logger.String("userId", userId.String()), logger.Int("numOrders", len(ordersToCancel)))
	return orderIds, nil
}
