package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// Generic Building blocks with no biz logic in a single TX
func (r *redisRepository) TxStart(ctx context.Context) (uint, error) {
	// --- START TRANSACTION ---
	tx := r.client.TxPipeline()
	r.ixIndex += 1
	txid := r.ixIndex
	r.txMap[txid] = tx

	return txid, nil
}

func (r *redisRepository) TxEnd(ctx context.Context, txid uint) {
	// --- END TRANSACTION ---
	if tx, ok := r.txMap[txid]; ok {
		_, err := tx.Exec(ctx)
		if err != nil {
			logctx.Error(ctx, "TxEnd exec failed", logger.Int("txid", int(txid)), logger.Error(err))
		}
		return
	}
	logctx.Error(ctx, "TxEnd txid not found", logger.Int("txid", int(txid)))
}
func (r *redisRepository) TxRemoveOrderFromPrice(ctx context.Context, txid uint, order models.Order) error {
	if tx, ok := r.txMap[txid]; ok {
		if order.Side == models.BUY {
			buyPricesKey := CreateBuySidePricesKey(order.Symbol)
			tx.ZRem(ctx, buyPricesKey, order.Id.String())
		} else {
			sellPricesKey := CreateSellSidePricesKey(order.Symbol)
			tx.ZRem(ctx, sellPricesKey, order.Id.String())
		}
		return nil
	} else {
		logctx.Error(ctx, "TxRemoveOrderFromPrice txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}
}
func (r *redisRepository) TxDeleteOrder(ctx context.Context, txid uint, orderId uuid.UUID) error {
	if tx, ok := r.txMap[txid]; ok {
		userOrdersKey := CreateUserOpenOrdersKey(orderId)
		tx.ZRem(ctx, userOrdersKey, orderId)
		return nil
	} else {
		logctx.Error(ctx, "TxDeleteOrder txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}
}

// stores only the state of the order
// as opposed to storeOrderTX - which should be deprecated in the end
func (r *redisRepository) TxStoreOrder(ctx context.Context, txid uint, order models.Order) error {
	if tx, ok := r.txMap[txid]; ok {
		// Store order details by order ID
		orderIDKey := CreateOrderIDKey(order.Id)
		orderMap := order.OrderToMap()
		tx.HSet(ctx, orderIDKey, orderMap)
		return nil
	} else {
		logctx.Error(ctx, "TxStoreOrder txid not found", logger.Int("txid", int(txid)))
		return models.ErrNotFound
	}
}
