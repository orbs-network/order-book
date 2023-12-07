package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	userOrdersKey := CreateUserOrdersKey(userId)

	// Fetch all order IDs for the user
	orderIdStrs, err := r.client.ZRange(ctx, userOrdersKey, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get order IDs for user", logger.String("userId", userId.String()), logger.Error(err))
		return nil, fmt.Errorf("failed to get order IDs for user: %v", err)
	}

	if len(orderIdStrs) == 0 {
		logctx.Warn(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return nil, models.ErrNoOrdersFound
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

		orders, err := r.FindOrdersByIds(ctx, orderIds[i:end])
		if err != nil {
			logctx.Error(ctx, "failed to find orders by IDs", logger.String("userId", userId.String()), logger.Error(err))
			return nil, fmt.Errorf("failed to find orders by IDs: %v", err)
		}

		ordersToCancel = append(ordersToCancel, orders...)
	}

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// Iterate through each order and remove associated data
	for _, order := range ordersToCancel {

		// Remove client order ID
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
		orderIDKey := CreateOrderIDKey(order.Id)
		transaction.Del(ctx, orderIDKey)

		logctx.Info(ctx, "removed order", logger.String("orderId", order.Id.String()), logger.String("symbol", order.Symbol.String()), logger.String("side", order.Side.String()), logger.String("userId", userId.String()))
	}

	// Remove the user's orders key
	transaction.Del(ctx, userOrdersKey)

	// Execute the transaction
	_, err = transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to remove user's orders in Redis", logger.String("userId", userId.String()), logger.Error(err))
		return nil, fmt.Errorf("failed to remove user's orders in Redis. Reason: %v", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "removed all orders for user", logger.String("userId", userId.String()), logger.Int("numOrders", len(ordersToCancel)))
	return orderIds, nil
}
