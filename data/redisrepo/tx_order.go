package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// Generic Building blocks with no biz logic in a single TX

// Perform a transaction with a single action. This should be used for all interactions with the Redis repository.
// Handles the transaction lifecycle.
// The action function should be a single Redis command or a series of Redis commands that should be executed in a single transaction.
// See the methods below (eg. TxModifyOrder, TxModifyPrices, etc.)
func (r *redisRepository) PerformTx(ctx context.Context, action func(txid uint) error) error {
	txid, err := r.txStart(ctx)
	if err != nil {
		logctx.Error(ctx, "PerformTransaction txStart failed", logger.Error(err))
		return fmt.Errorf("PerformTransaction txStart failed: %w", err)
	}
	defer func() {
		logctx.Debug(ctx, "PerformTransaction defer txEnd", logger.Int("txid", int(txid)))
		r.txEnd(ctx, txid)
	}()

	err = action(txid)
	if err != nil {
		logctx.Error(ctx, "PerformTransaction action failed", logger.Error(err), logger.Int("txid", int(txid)))
		return fmt.Errorf("PerformTransaction action failed: %w", err)
	}

	logctx.Debug(ctx, "PerformTransaction success", logger.Int("txid", int(txid)))
	return nil
}

// This should be used for all interactions with the `order:<id>` hash key
func (r *redisRepository) TxModifyOrder(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyOrder txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	switch operation {
	case models.Add, models.Update:
		// Store order details by order ID
		orderIDKey := CreateOrderIDKey(order.Id)
		orderMap := order.OrderToMap()
		tx.HSet(ctx, orderIDKey, orderMap)
		logctx.Debug(ctx, "TxModifyOrder add", logger.String("orderId", order.Id.String()), logger.String("orderMap", fmt.Sprintf("%v", orderMap)))
	case models.Remove:
		orderIDKey := CreateOrderIDKey(order.Id)
		tx.Del(ctx, orderIDKey)
		logctx.Debug(ctx, "TxModifyOrder remove", logger.String("orderId", order.Id.String()))
	default:
		logctx.Error(ctx, "TxModifyOrder unsupported operation", logger.Int("operation", int(operation)))
		return models.ErrUnsupportedOperation
	}

	return nil

}

// This should be used for all interactions with the `prices:<symbol>:buy` and `prices:<symbol>:sell` sorted sets (used to store bid/ask prices for each token pair)
func (r *redisRepository) TxModifyPrices(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyPrices txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	switch operation {
	case models.Add:
		// Add order to the sorted set for that token pair
		f64Price, _ := order.Price.Float64()
		timestamp := float64(order.Timestamp.UTC().UnixNano()) / 1e9
		score := f64Price + (timestamp / 1e12) // Use a combination of price and scaled timestamp so that orders with the same price are sorted by time. THIS SHOULD NOT BE USED FOR PRICE COMPARISON. Rather, use the price field in the order struct.

		if order.Side == models.BUY {
			buyPricesKey := CreateBuySidePricesKey(order.Symbol)
			tx.ZAdd(ctx, buyPricesKey, redis.Z{
				Score:  score,
				Member: order.Id.String(),
			})
		} else {
			sellPricesKey := CreateSellSidePricesKey(order.Symbol)
			tx.ZAdd(ctx, sellPricesKey, redis.Z{
				Score:  score,
				Member: order.Id.String(),
			})
		}
		logctx.Debug(ctx, "TxModifyPrices add", logger.String("orderId", order.Id.String()), logger.String("symbol", order.Symbol.String()), logger.String("side", order.Side.String()))
	case models.Remove:
		if order.Side == models.BUY {
			buyPricesKey := CreateBuySidePricesKey(order.Symbol)
			tx.ZRem(ctx, buyPricesKey, order.Id.String())
		} else {
			sellPricesKey := CreateSellSidePricesKey(order.Symbol)
			tx.ZRem(ctx, sellPricesKey, order.Id.String())
		}
		logctx.Debug(ctx, "TxModifyPrices remove", logger.String("orderId", order.Id.String()), logger.String("symbol", order.Symbol.String()), logger.String("side", order.Side.String()))
	default:
		logctx.Error(ctx, "TxModifyPrices unsupported operation", logger.Int("operation", int(operation)))
		return models.ErrUnsupportedOperation
	}
	return nil

}

// This should be used for all interactions with the `clientOID:<clientOID>` hash key
func (r *redisRepository) TxModifyClientOId(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyClientOId txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	switch operation {
	case models.Add:
		clientOIdKey := CreateClientOIDKey(order.ClientOId)
		tx.Set(ctx, clientOIdKey, order.Id.String(), 0)
		logctx.Debug(ctx, "ModifyClientOId add", logger.String("clientOID", order.ClientOId.String()), logger.String("orderId", order.Id.String()))
	case models.Remove:
		clientOIdKey := CreateClientOIDKey(order.ClientOId)
		tx.Del(ctx, clientOIdKey)
		logctx.Debug(ctx, "ModifyClientOId remove", logger.String("clientOID", order.ClientOId.String()), logger.String("orderId", order.Id.String()))
	default:
		logctx.Error(ctx, "ModifyClientOId unsupported operation", logger.Int("operation", int(operation)))
		return models.ErrUnsupportedOperation
	}
	return nil
}

// This should be used for all interactions with the `user:<userId>:openOrders` sorted set (used to store open orders for each user)
func (r *redisRepository) TxModifyUserOpenOrders(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyUserOpenOrders txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	switch operation {
	case models.Add:
		userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
		userOrdersScore := float64(order.Timestamp.UTC().UnixNano())
		tx.ZAdd(ctx, userOrdersKey, redis.Z{
			Score:  userOrdersScore,
			Member: order.Id.String(),
		})
		logctx.Debug(ctx, "ModifyUserOpenOrders add", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()))
	case models.Remove:
		userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
		tx.ZRem(ctx, userOrdersKey, order.Id.String())
		logctx.Debug(ctx, "ModifyUserOpenOrders remove", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()))
	default:
		logctx.Error(ctx, "ModifyUserOpenOrders unsupported operation", logger.Int("operation", int(operation)))
		return models.ErrUnsupportedOperation
	}
	return nil
}

// This should be used for all interactions with the `user:<userId>:filledOrders` sorted set (used to store partial-filled and cancelled OR fully filled orders for each user)
func (r *redisRepository) TxModifyUserFilledOrders(ctx context.Context, txid uint, operation models.Operation, order models.Order) error {
	var tx redis.Pipeliner
	var ok bool
	if tx, ok = r.txMap[txid]; !ok {
		logctx.Error(ctx, "TxModifyUserFilledOrders txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}

	switch operation {
	case models.Add:
		userFilledOrdersKey := CreateUserFilledOrdersKey(order.UserId)
		userFilledOrdersScore := float64(order.Timestamp.UTC().UnixNano())
		tx.ZAdd(ctx, userFilledOrdersKey, redis.Z{
			Score:  userFilledOrdersScore,
			Member: order.Id.String(),
		})
		logctx.Debug(ctx, "ModifyUserFilledOrders add", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()))
	case models.Remove:
		userFilledOrdersKey := CreateUserFilledOrdersKey(order.UserId)
		tx.ZRem(ctx, userFilledOrdersKey, order.Id.String())
		logctx.Debug(ctx, "ModifyUserFilledOrders remove", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()))
	default:
		logctx.Error(ctx, "ModifyUserFilledOrders unsupported operation", logger.Int("operation", int(operation)))
		return models.ErrUnsupportedOperation
	}
	return nil
}

// Removed an unfilled order
func (r *redisRepository) TxRemoveUnfilledOrder(ctx context.Context, txid uint, order models.Order) error {

	if err := r.TxModifyPrices(ctx, txid, models.Remove, order); err != nil {
		return err
	}

	if err := r.TxModifyUserOpenOrders(ctx, txid, models.Remove, order); err != nil {
		return err
	}

	if err := r.TxModifyClientOId(ctx, txid, models.Remove, order); err != nil {
		return err
	}

	if err := r.TxModifyOrder(ctx, txid, models.Remove, order); err != nil {
		return err
	}

	logctx.Debug(ctx, "TxRemoveUnfilledOrder", logger.String("orderId", order.Id.String()))
	return nil
}

func (r *redisRepository) txStart(ctx context.Context) (uint, error) {
	tx := r.client.TxPipeline()
	r.ixIndex += 1
	txid := r.ixIndex
	r.txMap[txid] = tx

	return txid, nil
}

func (r *redisRepository) txEnd(ctx context.Context, txid uint) {
	if tx, ok := r.txMap[txid]; ok {
		_, err := tx.Exec(ctx)
		if err != nil {
			logctx.Error(ctx, "txEnd exec failed", logger.Int("txid", int(txid)), logger.Error(err))
		}
		return
	}
	logctx.Error(ctx, "txEnd txid not found", logger.Int("txid", int(txid)))
}
