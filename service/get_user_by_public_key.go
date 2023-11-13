package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetUserByPublicKey(ctx context.Context, publicKey string) (*models.User, error) {

	user, err := s.orderBookStore.GetUserByPublicKey(ctx, publicKey)

	if err != nil {
		if err == models.ErrUserNotFound {
			logctx.Warn(ctx, "user not found", logger.String("publicKey", publicKey))
			return nil, err
		}

		logctx.Error(ctx, "unexpected error getting user by public key", logger.Error(err), logger.String("publicKey", publicKey))
		return nil, fmt.Errorf("unexpected error getting user by public key: %w", err)
	}

	if user == nil {
		logctx.Error(ctx, "user is nil but no error returned", logger.String("publicKey", publicKey))
		return nil, fmt.Errorf("user is nil but no error returned")
	}

	logctx.Info(ctx, "found user by public key", logger.String("userId", user.Id.String()))
	return user, nil
}
