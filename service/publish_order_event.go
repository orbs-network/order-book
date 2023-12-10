package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) publishOrderEvent(ctx context.Context, order *models.Order, status models.Status) {
	orderJson, err := order.ToJsonWithStatus(status)
	if err != nil {
		logctx.Error(ctx, "failed to marshal order to json", logger.Error(err))
	}

	if err := s.orderBookStore.PublishEvent(ctx, fmt.Sprintf("user_orders:%s", order.UserId), orderJson); err != nil {
		logctx.Error(ctx, "failed to publish order event", logger.Error(err))
	}
}
