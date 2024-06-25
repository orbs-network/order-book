package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// These methods should be used to store UNFILLED or PARTIALLY FILLED orders in Redis.
//
// `StoreFilledOrders` should be used to store completely filled orders.
func (r *redisRepository) txEnsureMakerTokenForBalanceTracking(ctx context.Context, txid uint, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyOrder txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	key := Order2MakerTokenTrackKey(order)
	if key == "" {
		logctx.Error(ctx, "Order2MakerTokenTrackKey failed for order", logger.String("orderId", order.Id.String()))
		return models.ErrInvalidInput
	}
	// Use SETNX to set the key only if it does not already exist
	result, err := tx.SetNX(ctx, key, -1, 0).Result()
	if err != nil {
		logctx.Error(ctx, "ensureMakerTokenForBalanceTracking failed to set key", logger.String("key", key), logger.Error(err))
		return err
	}

	if result {
		logctx.Debug(ctx, "MakerTokenTrackKey for balance was created value -1", logger.String("key", key))
	}

	return nil
}

func (r *redisRepository) txAddOpenOrder(ctx context.Context, txid uint, order models.Order) error {
	// add to order:id key
	if err := r.TxModifyOrder(ctx, txid, models.Add, order); err != nil {
		logctx.Error(ctx, "StoreOpenOrders TxModifyOrder Failed adding order", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
		return err
	}
	// add to client oid key
	if err := r.TxModifyClientOId(ctx, txid, models.Add, order); err != nil {
		logctx.Error(ctx, "StoreOpenOrders TxModifyClientOId Failed adding order", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
		return err
	}

	// add to price
	if err := r.TxModifyPrices(ctx, txid, models.Add, order); err != nil {
		logctx.Error(ctx, "StoreOpenOrders TxModifyPrices Failed adding order", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
		return err
	}
	// add to user:opeOrder key
	if err := r.TxModifyUserOpenOrders(ctx, txid, models.Add, order); err != nil {
		logctx.Error(ctx, "TxModifyUserOpenOrders TxModifyClientOId Failed adding order", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
		return err
	}

	// ensure balance is tracked for this order
	return r.txEnsureMakerTokenForBalanceTracking(ctx, txid, order)
}

func (r *redisRepository) StoreOpenOrders(ctx context.Context, orders []models.Order) error {
	err := r.PerformTx(ctx, func(txid uint) error {
		for _, order := range orders {
			if err := r.txAddOpenOrder(ctx, txid, order); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (r *redisRepository) StoreOpenOrder(ctx context.Context, order models.Order) error {
	err := r.PerformTx(ctx, func(txid uint) error {
		return r.txAddOpenOrder(ctx, txid, order)
	})

	return err
}
