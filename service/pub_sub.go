package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) SubscribeUserOrders(ctx context.Context, userId uuid.UUID) (chan []byte, error) {
	logctx.Debug(ctx, "subscribing to user orders", logger.String("userId", userId.String()))

	eventKey := models.CreateUserOrdersEventKey(userId)

	channel, err := s.orderBookStore.SubscribeToEvents(ctx, fmt.Sprintf("user_orders:%s", userId))
	if err != nil {
		logctx.Error(ctx, "failed to subscribe to user orders", logger.String("event", eventKey), logger.Error(err))
		return nil, fmt.Errorf("failed to subscribe to user orders: %w", err)
	}

	return channel, nil
}

func (s *EvmClient) publishOrderEvent(ctx context.Context, order *models.Order) {
	key, value, err := createOrderEvent(ctx, order)
	if err != nil {
		return
	}

	if err := s.orderBookStore.PublishEvent(ctx, key, value); err != nil {
		logctx.Error(ctx, "failed to publish order event", logger.String("event", key), logger.Error(err))
	}
}

func (s *Service) publishOrderEvent(ctx context.Context, order *models.Order) {
	key, value, err := createOrderEvent(ctx, order)
	if err != nil {
		return
	}

	if err := s.orderBookStore.PublishEvent(ctx, key, value); err != nil {
		logctx.Error(ctx, "failed to publish order event", logger.String("event", key), logger.Error(err))
	}
}

func createOrderEvent(ctx context.Context, order *models.Order) (key string, value []byte, err error) {
	value, err = order.ToJson()
	if err != nil {
		logctx.Error(ctx, "failed to marshal order to json", logger.Error(err))
	}

	key = models.CreateUserOrdersEventKey(order.UserId)

	return key, value, err
}
