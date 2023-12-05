package serviceuser

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetUserByApiKey(ctx context.Context, apiKey string) (*models.User, error) {

	user, err := s.userStore.GetUserByApiKey(ctx, apiKey)

	if err == models.ErrUserNotFound {
		logctx.Warn(ctx, "user not found", logger.Error(err))
		return nil, err
	}

	if err != nil {
		logctx.Error(ctx, "failed to get user by api key", logger.Error(err))
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}

	logctx.Info(ctx, "user retrieved by API key", logger.String("userId", user.Id.String()), logger.String("pubKey", user.PubKey))
	return user, nil
}
