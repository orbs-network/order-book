package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// Order remains in DB, but is marked as cancelled
// Order is removed from the prices sorted set, user's order set and order hash is updated to cancelled
// Upon swap resolve false -> should be removed
func (r *redisRepository) CancelPendingOrder(ctx context.Context, order models.Order) error {

	if !order.IsPending() {
		logctx.Error(ctx, "trying to cancel pending order which is not pending", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		return models.ErrOrderNotPending
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

	_, err := transaction.Exec(ctx)

	if err != nil {
		logctx.Error(ctx, "failed to cancel pending order", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("failed to cancel pending order: %w", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "cancelled pending order", logger.String("userId", order.UserId.String()), logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}
