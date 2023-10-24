package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

var ErrInvalidOrderSide = fmt.Errorf("invalid order side")

func (r *redisRepository) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (models.Order, error) {
	var key string
	var err error
	var orderIDs []string

	if side == models.BUY {
		key = CreateBuySidePricesKey(symbol)
	} else if side == models.SELL {
		key = CreateSellSidePricesKey(symbol)
	} else {
		return models.Order{}, ErrInvalidOrderSide
	}

	if side == models.BUY {
		// Highest bid price for buying
		orderIDs, err = r.client.ZRevRange(ctx, key, 0, 0).Result()
	} else {
		// Lowest ask price for selling
		orderIDs, err = r.client.ZRange(ctx, key, 0, 0).Result()
	}

	if err != nil {
		return models.Order{}, err
	}

	if len(orderIDs) == 0 {
		return models.Order{}, models.ErrOrderNotFound
	}

	orderIDStr := orderIDs[0]

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return models.Order{}, err
	}

	order, err := r.FindOrderById(ctx, orderID, false)

	if err != nil {
		logctx.Error(ctx, "unexpected error when getting order for best price", logger.String("orderId", orderID.String()), logger.Error(err))
		return models.Order{}, err
	}

	return *order, nil
}
