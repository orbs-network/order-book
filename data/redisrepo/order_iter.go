package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

//////////////////////////////////////////////////
type OrderIter struct {
	index uint
	ids   []string
	redis *redisRepository
}

func (i OrderIter) Next() *models.Order {
	ctx := context.Background()

	if i.index >= uint(len(i.ids)) {
		logctx.Error(ctx, "Error iterator reached last element")
		return nil
	}

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
	// increment index
	i.index += 1

	return order
}

func (i OrderIter) HasNext() bool {
	return i.index < uint(len(i.ids))-1
}

//////////////////////////////////////////////////
func (r *redisRepository) GetMinAsk(ctx context.Context, symbol models.Symbol) service.OrderIter {

	key := CreateSellSidePricesKey(symbol)

	// Min ask price for selling
	orderIDs, err := r.client.ZRange(ctx, key, 0, 0).Result()
	if err != nil {
		logctx.Error(ctx, "Error fetching asks", logger.Error(err))
	}
	// create order iter
	return &OrderIter{
		index: 0,
		ids:   orderIDs,
		redis: r,
	}

}

//////////////////////////////////////////////////
func (r *redisRepository) GetMaxBid(ctx context.Context, symbol models.Symbol) service.OrderIter {
	key := CreateBuySidePricesKey(symbol)

	// Min ask price for selling
	orderIDs, err := r.client.ZRevRange(ctx, key, 0, 0).Result()

	if err != nil {
		logctx.Error(ctx, "Error fetching bids", logger.Error(err))
	}
	// create order iter
	return &OrderIter{
		ids:   orderIDs,
		redis: r,
	}
}
