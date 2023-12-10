package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) SubscribeUserOrders(ctx context.Context, userId uuid.UUID) (chan []byte, error) {
	logctx.Info(ctx, "subscribing to user orders", logger.String("userId", userId.String()))

	channel, err := s.orderBookStore.SubscribeToEvents(ctx, fmt.Sprintf("user_orders:%s", userId))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to user orders: %w", err)
	}

	return channel, nil
}
