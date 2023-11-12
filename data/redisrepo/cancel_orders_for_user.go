package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// TODO: this is a fairly expensive operation. We should use bulk operations to improve performance, and consider batching removals in case a of a large number of orders.
func (r *redisRepository) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) error {
	userOrdersKey := CreateUserOrdersKey(userId)

	// Fetch all order IDs for the user
	orderIdStrs, err := r.client.ZRange(ctx, userOrdersKey, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get order IDs for user", logger.String("userId", userId.String()), logger.Error(err))
		return fmt.Errorf("failed to fetch user order IDs. Reason: %v", err)
	}

	if len(orderIdStrs) == 0 {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return nil
	}

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// Iterate through each order ID and remove associated data
	for _, orderIdStr := range orderIdStrs {
		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "failed to parse order ID", logger.String("orderId", orderIdStr), logger.Error(err))
			return fmt.Errorf("failed to parse order ID: %v", err)
		}

		// Remove client order ID
		order, err := r.FindOrderById(ctx, orderId, false)
		if err == models.ErrOrderNotFound {
			logctx.Error(ctx, "order should exist but could not be found by ID. Continuing with user's other orders", logger.String("orderId", orderIdStr), logger.Error(err))
			continue
		}
		if err != nil {
			logctx.Error(ctx, "unexpected error finding order by ID", logger.String("orderId", orderIdStr), logger.Error(err))
			return fmt.Errorf("unexpected error finding order by ID: %v", err)
		}

		clientOIDKey := CreateClientOIDKey(order.ClientOId)
		transaction.Del(ctx, clientOIDKey)

		if order.Side == models.BUY {
			buyPricesKey := CreateBuySidePricesKey(order.Symbol)
			transaction.ZRem(ctx, buyPricesKey, order.Id.String())
		} else {
			sellPricesKey := CreateSellSidePricesKey(order.Symbol)
			transaction.ZRem(ctx, sellPricesKey, order.Id.String())
		}

		// Remove order details by order ID
		orderIDKey := CreateOrderIDKey(orderId)
		transaction.Del(ctx, orderIDKey)

		logctx.Info(ctx, "removed order", logger.String("orderId", orderIdStr), logger.String("symbol", order.Symbol.String()), logger.String("side", order.Side.String()), logger.String("userId", userId.String()))
	}

	// Remove the user's orders key
	transaction.Del(ctx, userOrdersKey)

	// Execute the transaction
	_, err = transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to remove user's orders in Redis", logger.String("userId", userId.String()), logger.Error(err))
		return fmt.Errorf("failed to remove user's orders in Redis. Reason: %v", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "removed all orders for user", logger.String("userId", userId.String()))
	return nil
}
