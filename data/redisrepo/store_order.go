package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

func (r *redisRepository) StoreOrder(ctx context.Context, order models.Order) error {

	orderMap := order.OrderToMap()

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()

	// Keep track of that user's orders
	userOrdersKey := CreateUserOrdersKey(order.UserId)
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

	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to store order in Redis", logger.Error(err), logger.String("orderId", order.Id.String()))
		return fmt.Errorf("transaction failed. Reason: %v", err)
	}
	// --- END TRANSACTION ---

	logctx.Info(ctx, "stored order", logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}
