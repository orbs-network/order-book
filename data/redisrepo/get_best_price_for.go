package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

func (r *redisRepository) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error) {
	var key string
	var price float64
	var err error
	var prices []redis.Z

	if side == models.BUY {
		key = CreateBuySidePricesKey(symbol)
	} else if side == models.SELL {
		key = CreateSellSidePricesKey(symbol)
	} else {
		return decimal.Zero, fmt.Errorf("invalid order side")
	}

	if side == models.BUY {
		// Highest bid price for buying
		prices, err = r.client.ZRevRangeWithScores(ctx, key, 0, 0).Result()
	} else {
		// Lowest ask price for selling
		prices, err = r.client.ZRangeWithScores(ctx, key, 0, 0).Result()
	}

	if err != nil {
		return decimal.Zero, err
	}

	if len(prices) == 0 {
		return decimal.Zero, models.ErrOrderNotFound
	}

	price = prices[0].Score

	decimalPrice := decimal.NewFromFloat(price)
	return decimalPrice, nil
}
