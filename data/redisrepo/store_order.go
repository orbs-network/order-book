package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

func (r *redisRepository) StoreOrder(ctx context.Context, order models.Order) error {

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("could not marshal order: %v", err)
	}

	// --- Start transaction ---
	pipeline := r.client.TxPipeline()

	// Price hash
	priceKey := CreatePriceKey(order.Symbol, order.Price)
	pipeline.HSet(ctx, priceKey, order.Id, orderJSON)

	// Order ID hash
	orderIDKey := CreateOrderIDKey(order.Symbol)
	pipeline.HSet(ctx, orderIDKey, order.Id, orderJSON)

	// Prices sorted set
	r.updatePricesSortedSet(order)

	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}
	// --- End transaction ---

	logctx.Info(ctx, "stored order", logger.String("order", string(orderJSON)))
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
