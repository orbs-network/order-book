package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrdersForUser(ctx context.Context, publicKey string) error {

	user, err := s.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		logctx.Warn(ctx, "user not found", logger.String("publicKey", publicKey), logger.Error(err))
		return err
	}

	if err = s.orderBookStore.CancelOrdersForUser(ctx, user.Id); err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		return fmt.Errorf("could not cancel orders for user: %w", err)
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", user.Id.String()))
	return nil
}
