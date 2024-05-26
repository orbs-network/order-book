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
	// continue iterating even if error occured, next order may be healthy
	for i.index < len(i.ids)-1 {
		// increment index - first one is -1
		i.index = i.index + 1
		// get order
		orderId, err := uuid.Parse(i.ids[i.index])
		if err != nil {
			logctx.Error(ctx, "Error parsing order id", logger.Error(err))
		} else {
			order, err := i.redis.FindOrderById(ctx, orderId, false)
			if err != nil {
				logctx.Warn(ctx, "Order not found, perhaps deleted, go next", logger.String("orderId", orderId.String()), logger.Error(err))
			} else {
				// success
				return order
			}
		}
	}
	logctx.Warn(ctx, "Error iterator reached last element")
	return nil
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
		return nil
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
		return nil
	}
	// create order iter
	return &OrderIter{
		index: -1,
		ids:   orderIDs,
		redis: r,
	}
}
