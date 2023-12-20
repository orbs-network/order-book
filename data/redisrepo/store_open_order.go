package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// These methods should be used to store UNFILLED or PARTIALLY FILLED orders in Redis.
//
// `StoreFilledOrders` should be used to store completely filled orders.
//
// TODO: combine `StoreOpenOrder` and `StoreFilledOrder` into a single `StoreOrder` method that checks order status and stores accordingly.
func (r *redisRepository) StoreOpenOrder(ctx context.Context, order models.Order) error {

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	err := storeOrderTX(ctx, transaction, &order)
	if err != nil {
		return err
	}

	_, err = transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to store open order in Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("failed to stores open order in Redis: %v", err)
	}
	logctx.Info(ctx, "stored order", logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}

func (r *redisRepository) StoreOpenOrders(ctx context.Context, orders []models.Order) error {
	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	for _, order := range orders {
		err := storeOrderTX(ctx, transaction, &order)
		if err != nil {
			return err
		}

		logctx.Info(ctx, "stored order", logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))

	}
	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to stores open order in Redis", logger.Error(err), logger.Strings("orderIds", models.OrderIdsToStrings(ctx, &orders)))
		return fmt.Errorf("failed to stores open order in Redis: %v", err)
	}
	return nil
}

func storeOrderTX(ctx context.Context, transaction redis.Pipeliner, order *models.Order) error {
	orderMap := order.OrderToMap()

	// Keep track of that user's orders
	userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
	userOrdersScore := float64(order.Timestamp.UTC().UnixNano())
	transaction.ZAdd(ctx, userOrdersKey, redis.Z{
		Score:  userOrdersScore,
		Member: order.Id.String(),
	})

	// Store order details by order ID
	orderIDKey := CreateOrderIDKey(order.Id)
	transaction.HSet(ctx, orderIDKey, orderMap)

	// Store client order ID
	clientOIDKey := CreateClientOIDKey(order.ClientOId)
	transaction.Set(ctx, clientOIDKey, order.Id.String(), 0)

	// Add order to the sorted set for that token pair
	f64Price, _ := order.Price.Float64()
	timestamp := float64(order.Timestamp.UTC().UnixNano()) / 1e9
	score := f64Price + (timestamp / 1e12) // Use a combination of price and scaled timestamp so that orders with the same price are sorted by time. This should not be used for price comparison.

	if order.Side == models.BUY {
		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		transaction.ZAdd(ctx, buyPricesKey, redis.Z{
			Score:  score,
			Member: order.Id.String(),
		})
	} else {
		sellPricesKey := CreateSellSidePricesKey(order.Symbol)
		transaction.ZAdd(ctx, sellPricesKey, redis.Z{
			Score:  score,
			Member: order.Id.String(),
		})
	}

	return nil
}
