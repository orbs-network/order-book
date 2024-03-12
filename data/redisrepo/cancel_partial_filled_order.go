package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// Cancels a partial filled order.
// Order is removed from the prices sorted set, user's order set and order hash is updated to cancelled
// can be called to locked orders
func (r *redisRepository) CancelPartialFilledOrder(ctx context.Context, order models.Order) error {

	if order.IsPending() {
		logctx.Error(ctx, "trying to cancel order that is currently pending", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		return models.ErrOrderPending
	}

	if !order.IsPartialFilled() {
		logctx.Error(ctx, "trying to cancel order that is not partial filled", logger.String("orderId", order.Id.String()), logger.String("sizeFilled", order.SizeFilled.String()))
		return models.ErrOrderNotPartialFilled
	}

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// update order hash
	orderIDKey := CreateOrderIDKey(order.Id)
	transaction.HSet(ctx, orderIDKey, "cancelled", "true")

	// remove from price set
	if order.Side == models.BUY {
		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		transaction.ZRem(ctx, buyPricesKey, order.Id.String())
	} else {
		sellPricesKey := CreateSellSidePricesKey(order.Symbol)
		transaction.ZRem(ctx, sellPricesKey, order.Id.String())
	}

	// remove from user order set
	userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
	transaction.ZRem(ctx, userOrdersKey, order.Id.String())

	// add to filled orders set
	userFilledOrdersKey := CreateUserFilledOrdersKey(order.UserId)
	userFilledOrdersScore := float64(order.Timestamp.UTC().UnixNano())
	transaction.ZAdd(ctx, userFilledOrdersKey, redis.Z{
		Score:  userFilledOrdersScore,
		Member: order.Id.String(),
	})

	_, err := transaction.Exec(ctx)

	if err != nil {
		logctx.Error(ctx, "failed to cancel partial filled order", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("failed to cancel partial filled order: %w", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "cancelled partial filled order", logger.String("userId", order.UserId.String()), logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}
