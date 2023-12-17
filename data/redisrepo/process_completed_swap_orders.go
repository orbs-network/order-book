package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// ProcessCompletedSwapOrders stores the updated swap orders and removes the swap from Redis. It should be called after a swap is completed.
//
// `orders` should be the orders that were part of the swap (with `SizePending` and `SizeFilled` updated accordingly)
//
// `isSuccessful` should be `true` if the swap was successful, `false` otherwise
func (r *redisRepository) ProcessCompletedSwapOrders(ctx context.Context, orders []*models.Order, swapId uuid.UUID, isSuccessful bool) error {
	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// 1. Store updated swap orders, handle completely filled orders
	if isSuccessful {
		for _, order := range orders {
			if order.IsFilled() {
				if err := storeFilledOrderTx(ctx, transaction, order); err != nil {
					logctx.Error(ctx, "failed to store filled order in Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
					return fmt.Errorf("failed to store filled order in Redis: %v", err)
				}
			} else {
				if err := storeOrderTX(ctx, transaction, order); err != nil {
					logctx.Error(ctx, "failed to store open order in Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
					return fmt.Errorf("failed to store open order in Redis: %v", err)
				}
			}
		}
	} else {
		// Store updated orders
		for _, order := range orders {
			if err := storeOrderTX(ctx, transaction, order); err != nil {
				logctx.Error(ctx, "failed to store open order in Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
				return fmt.Errorf("failed to store open order in Redis: %v", err)
			}
		}
	}

	// 2. Remove the swap
	swapKey := CreateSwapKey(swapId)
	transaction.Del(ctx, swapKey).Err()

	// --- END TRANSACTION ---
	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute ProcessCompletedSwapOrders transaction", logger.Error(err))
		return fmt.Errorf("failed to execute ProcessCompletedSwapOrders transaction: %v", err)
	}

	return nil
}
