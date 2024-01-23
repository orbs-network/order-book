package serviceuser

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CreateUserInput struct {
	PubKey string
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (user models.User, apiKey string, err error) {

	apiKey, err = GenerateAPIKey()
	if err != nil {
		logctx.Error(ctx, "failed to generate api key", logger.Error(err))
		return models.User{}, "", fmt.Errorf("failed to generate api key: %w", err)
	}

	userId := uuid.New()

	createdUser, err := s.userStore.CreateUser(ctx, models.User{
		Id:     userId,
		PubKey: input.PubKey,
		Type:   models.MARKET_MAKER,
		ApiKey: HashAPIKey(apiKey),
	})

	if err == models.ErrUserAlreadyExists {
		logctx.Warn(ctx, "user already exists for that pub key", logger.String("pubKey", input.PubKey))
		return models.User{}, "", err
	}

	if err != nil {
		logctx.Error(ctx, "failed to create user", logger.String("userId", userId.String()), logger.Error(err))
		return models.User{}, "", fmt.Errorf("failed to create user: %w", err)
	}

	logctx.Info(ctx, "user created", logger.String("userId", createdUser.Id.String()), logger.String("pubKey",
		createdUser.PubKey))

	return createdUser, apiKey, nil
}
