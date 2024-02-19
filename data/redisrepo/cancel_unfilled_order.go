package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// Cancels an unfilled order.
// Order is removed from the prices sorted set, user's order set and order hash is removed
// May only be called if order is not pending and completely unfilled
func (r *redisRepository) CancelUnfilledOrder(ctx context.Context, order models.Order) error {
	if order.IsPending() {
		logctx.Error(ctx, "trying to remove order that is currently pending", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		return models.ErrOrderPending
	}

	if !order.IsUnfilled() {
		logctx.Error(ctx, "trying to remove order that is not unfilled", logger.String("orderId", order.Id.String()), logger.String("sizeFilled", order.SizeFilled.String()))
		return models.ErrOrderNotUnfilled
	}

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()
	if order.Side == models.BUY {
		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		transaction.ZRem(ctx, buyPricesKey, order.Id.String())
	} else {
		sellPricesKey := CreateSellSidePricesKey(order.Symbol)
		transaction.ZRem(ctx, sellPricesKey, order.Id.String())
	}

	userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
	transaction.ZRem(ctx, userOrdersKey, order.Id.String())

	clientOIdKey := CreateClientOIDKey(order.ClientOId)
	transaction.Del(ctx, clientOIdKey)

	// remove order hash
	orderIDKey := CreateOrderIDKey(order.Id)
	transaction.Del(ctx, orderIDKey)

	_, err := transaction.Exec(ctx)

	if err != nil {
		logctx.Error(ctx, "failed to remove order from Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("failed to remove order from Redis: %w", err)
	}
	// --- END TRANSACTION ---

	logctx.Debug(ctx, "removed unfilled order", logger.String("userId", order.UserId.String()), logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}
