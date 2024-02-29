package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type ProcessCompletedSwapOrdersInput struct {
	IsSuccessful bool
	Orders       []*models.Order
	SwapId       uuid.UUID
	Tx           models.Tx
}

// ProcessCompletedSwapOrders stores the updated swap orders and removes the swap from Redis. It should be called after a swap is completed.
//
// `orders` should be the orders that were part of the swap (with `SizePending` and `SizeFilled` updated accordingly)
//
// `isSuccessful` should be `true` if the swap was successful, `false` otherwise
func (r *redisRepository) ProcessCompletedSwapOrders(ctx context.Context, ordersWithSize []store.OrderWithSize, swapId uuid.UUID, tx *models.Tx, isSuccessful bool) error {
	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// 1a. Store updated swap orders, handle completely filled orders
	if isSuccessful {
		for _, o := range ordersWithSize {
			if o.Order.IsFilled() {
				if err := storeFilledOrderTx(ctx, transaction, o.Order); err != nil {
					logctx.Error(ctx, "failed to store filled order in Redis", logger.Error(err), logger.String("orderId", o.Order.Id.String()))
					return fmt.Errorf("failed to store filled order in Redis: %v", err)
				}
			} else {
				if err := storeOrderTX(ctx, transaction, o.Order); err != nil {
					logctx.Error(ctx, "failed to store open order in Redis", logger.Error(err), logger.String("orderId", o.Order.Id.String()))
					return fmt.Errorf("failed to store open order in Redis: %v", err)
				}
			}
			// 1b. Store completed swap details
			err := r.StoreCompletedSwap(ctx, store.StoreCompletedSwapInput{
				UserId:    o.Order.UserId,
				SwapId:    swapId,
				OrderId:   o.Order.Id,
				TxId:      tx.TxHash,
				Timestamp: *tx.Timestamp,
				Block:     *tx.Block,
			})
			if err != nil {
				// TODO: should we return an error here?
				logctx.Error(ctx, "failed to store successful tx completed swap in Redis", logger.Error(err), logger.String("swapId", swapId.String()), logger.String("orderId", o.Order.Id.String()), logger.String("txHash", tx.TxHash))
			}
		}
	} else {
		// 1a. Store updated orders
		for _, o := range ordersWithSize {
			if err := storeOrderTX(ctx, transaction, o.Order); err != nil {
				logctx.Error(ctx, "failed to store open order in Redis", logger.Error(err), logger.String("orderId", o.Order.Id.String()))
				return fmt.Errorf("failed to store open order in Redis: %v", err)
			}
			// 1b. Store completed swap details
			err := r.StoreCompletedSwap(ctx, store.StoreCompletedSwapInput{
				UserId:    o.Order.UserId,
				SwapId:    swapId,
				OrderId:   o.Order.Id,
				Size:      o.Size,
				TxId:      tx.TxHash,
				Timestamp: *tx.Timestamp,
				Block:     *tx.Block,
			})
			if err != nil {
				// TODO: should we return an error here?
				logctx.Error(ctx, "failed to store failed tx completed swap in Redis", logger.Error(err), logger.String("swapId", swapId.String()), logger.String("orderId", o.Order.Id.String()), logger.String("txHash", tx.TxHash))
			}
		}
	}

	// 2. Remove the swap
	swapKey := CreateSwapKey(swapId)
	transaction.Del(ctx, swapKey)

	// --- END TRANSACTION ---
	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute ProcessCompletedSwapOrders transaction", logger.Error(err))
		return fmt.Errorf("failed to execute ProcessCompletedSwapOrders transaction: %v", err)
	}

	return nil
}

func (r *redisRepository) ResolveSwap(ctx context.Context, swap models.Swap) error {

	// save swap in resolved key
	err := r.saveSwap(ctx, swap.Id, swap, true)
	if err != nil {
		logctx.Error(ctx, "failed to save swap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return err
	}

	// remove from swapId
	err = r.RemoveSwap(ctx, swap.Id)
	if err != nil {
		logctx.Error(ctx, "failed to remove swap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return err
	}

	return nil
}

// save swapId in a set of the userId:resolvedSwap key
func (r *redisRepository) StoreUserResolvedSwap(ctx context.Context, userId uuid.UUID, swap models.Swap) error {
	key := CreateUserResolvedSwapsKey(userId)
	return AddVal2Set(ctx, r.client, key, swap.Id.String())
}
