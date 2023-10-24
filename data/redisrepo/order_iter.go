package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// ////////////////////////////////////////////////
type OrderIter struct {
	index int
	ids   []string
	redis *redisRepository
}

func (i *OrderIter) Next(ctx context.Context) *models.Order {
	//ctx := context.Background()

	if i.index >= len(i.ids) {
		logctx.Error(ctx, "Error iterator reached last element")
		return nil
	}

	// increment index
	i.index = i.index + 1
	// get order
	orderId, err := uuid.Parse(i.ids[i.index])
	if err != nil {
		logctx.Error(ctx, "Error parsing bid order id", logger.Error(err))
		return nil
	}
	order, err := i.redis.FindOrderById(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "Error fetching order", logger.Error(err))
		return nil
	}

	return order
}

func (i *OrderIter) HasNext() bool {
	return i.index < (len(i.ids) - 1)
}

// ////////////////////////////////////////////////
func (r *redisRepository) GetMinAsk(ctx context.Context, symbol models.Symbol) models.OrderIter {
	key := CreateSellSidePricesKey(symbol)

	// Min ask price for selling
	orderIDs, err := r.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "Error fetching asks", logger.Error(err))
	}
	// create order iter
	return &OrderIter{
		index: -1,
		ids:   orderIDs,
		redis: r,
	}

}

// ////////////////////////////////////////////////
func (r *redisRepository) GetMaxBid(ctx context.Context, symbol models.Symbol) models.OrderIter {
	key := CreateBuySidePricesKey(symbol)

	// Min ask price for selling
	orderIDs, err := r.client.ZRevRange(ctx, key, 0, -1).Result()

	if err != nil {
		logctx.Error(ctx, "Error fetching bids", logger.Error(err))
	}
	// create order iter
	return &OrderIter{
		index: -1,
		ids:   orderIDs,
		redis: r,
	}
}
