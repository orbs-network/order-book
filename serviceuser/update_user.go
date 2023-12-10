package serviceuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type UpdateUserInput struct {
	UserId uuid.UUID
	ApiKey string
	PubKey string
}

func (s *Service) UpdateUser(ctx context.Context, input UpdateUserInput) error {

	if input.ApiKey == "" || input.PubKey == "" {
		logctx.Error(ctx, "apiKey or pubKey is empty", logger.String("apiKey", input.ApiKey), logger.String("pubKey", input.PubKey))
		return models.ErrInvalidInput
	}

	if err := s.userStore.UpdateUser(ctx, storeuser.UpdateUserInput{
		UserId: input.UserId,
		ApiKey: input.ApiKey,
		PubKey: input.PubKey,
	}); err != nil {
		logctx.Error(ctx, "failed to update user", logger.Error(err), logger.String("userId", input.UserId.String()))
		return fmt.Errorf("failed to update user: %w", err)
	}

	logctx.Info(ctx, "user updated", logger.String("userId", input.UserId.String()), logger.String("pubKey", input.PubKey))
	return nil
}
