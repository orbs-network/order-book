package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// Removes an order from the order book.
// Order is removed from the prices sorted set, user's order set and order hash is updated to `status: CANCELED`
// SHOULD ONLY BE USED WHEN ORDER STATUS IS STILL `OPEN`
func (r *redisRepository) RemoveOrder(ctx context.Context, order models.Order) error {

	if order.Status != models.STATUS_OPEN {
		logctx.Error(ctx, "trying to remove order that is not open", logger.String("orderId", order.Id.String()), logger.String("status", order.Status.String()))
		return models.ErrOrderNotOpen
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

	userOrdersKey := CreateUserOrdersKey(order.UserId)
	transaction.ZRem(ctx, userOrdersKey, order.Id.String())

	clientOIdKey := CreateClientOIDKey(order.ClientOId)
	transaction.Del(ctx, clientOIdKey, order.ClientOId.String())

	// update order status to CANCELED
	orderIDKey := CreateOrderIDKey(order.Id)
	//remove order instead
	transaction.Del(ctx, orderIDKey)

	_, err := transaction.Exec(ctx)

	if err != nil {
		logctx.Error(ctx, "failed to remove order from Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("failed to remove order from Redis: %w", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "removed order", logger.String("userId", order.UserId.String()), logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}
