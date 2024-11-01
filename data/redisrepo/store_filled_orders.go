package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// This should be used to store FILLED orders in Redis.
//
// `StoreOpenOrder` or `StoreOpenOrders` should be used to store unfilled or partially filled orders.
func (r *redisRepository) StoreFilledOrders(ctx context.Context, orders []models.Order) error {
	transaction := r.client.TxPipeline()

	for _, order := range orders {
		err := storeFilledOrderTx(ctx, transaction, &order)
		if err != nil {
			return err
		}

	}
	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to store filled orders in Redis", logger.Error(err), logger.Strings("orderIds", models.OrderIdsToStrings(ctx, &orders)))
		return fmt.Errorf("failed to store filled orders in Redis: %v", err)
	}

	logctx.Debug(ctx, "stored filled orders in Redis", logger.Strings("orderIds", models.OrderIdsToStrings(ctx, &orders)))
	return nil
}

func storeFilledOrderTx(ctx context.Context, transaction redis.Pipeliner, order *models.Order) error {
	// 1. Remove the order from the user's open orders set
	userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
	transaction.ZRem(ctx, userOrdersKey, order.Id.String())
	// 2. Remove the order from the buy/sell prices set for that pair
	if order.Side == models.BUY {
		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		transaction.ZRem(ctx, buyPricesKey, order.Id.String())
	} else {
		sellPricesKey := CreateSellSidePricesKey(order.Symbol)
		transaction.ZRem(ctx, sellPricesKey, order.Id.String())
	}

	// 4. Store the order in the order ID key
	orderIDKey := CreateOrderIDKey(order.Id)
	transaction.HSet(ctx, orderIDKey, order.OrderToMap())

	return nil
}
