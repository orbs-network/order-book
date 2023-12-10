package serviceuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {

	user, err := s.userStore.GetUserById(ctx, userId)
	if err != nil {
		logctx.Error(ctx, "failed to get user by id", logger.Error(err), logger.String("userId", userId.String()))
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	logctx.Info(ctx, "user retrieved by ID", logger.String("userId", userId.String()), logger.String("pubKey", user.PubKey))
	return user, nil
}
