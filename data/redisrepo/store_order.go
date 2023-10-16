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

	// --- Start transaction ---
	transaction := r.client.TxPipeline()

	// Price hash
	priceKey := CreatePriceKey(order.Symbol, order.Price)
	for k, v := range orderMap {
		transaction.HSet(ctx, priceKey, k, v).Err()
	}

	// Order ID hash
	orderIDKey := CreateOrderIDKey(order.Id)
	for k, v := range orderMap {
		transaction.HSet(ctx, orderIDKey, k, v).Err()
	}

	// Prices sorted set
	r.updatePricesSortedSet(order)

	_, err := transaction.Exec(ctx)
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}
	// --- End transaction ---

	logctx.Info(ctx, "stored order", logger.String("orderId", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()), logger.String("side", order.Side.String()))
	return nil
}

func (r *redisRepository) updatePricesSortedSet(order models.Order) {

	float64Price, _ := order.Price.Float64()

	if order.Side == models.BUY {
		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		r.client.ZAdd(context.Background(), buyPricesKey, redis.Z{
			Score:  float64Price,
			Member: order.Id,
		})
	} else {
		sellPricesKey := CreateSellSidePricesKey(order.Symbol)
		r.client.ZAdd(context.Background(), sellPricesKey, redis.Z{
			Score:  float64Price,
			Member: order.Id,
		})
	}
}
