package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrdersForUser(ctx context.Context, userId uuid.UUID) error {

	if err := s.orderBookStore.CancelOrdersForUser(ctx, userId); err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", userId.String()))
		return fmt.Errorf("could not cancel orders for user: %w", err)
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", userId.String()))
	return nil
}
