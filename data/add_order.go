package data

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *memoryRespository) AddOrder(ctx context.Context, order models.Order) (models.Order, error) {
	orderIdStr := order.Id.String()
	logctx.Info(ctx, "adding order to in-memory DB", logger.String("orderId", orderIdStr))
	r.Orders[orderIdStr] = order
	return order, nil
}
