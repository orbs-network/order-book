package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// GetUserByApiKey returns a user by their apiKey
func (r *redisRepository) GetUserByApiKey(ctx context.Context, apiKey string) (*models.User, error) {

	key := CreateUserApiKeyKey(apiKey)

	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		logctx.Error(ctx, "unexpected error getting user by api key", logger.Error(err))
		return nil, fmt.Errorf("unexpected error getting user by api key: %w", err)
	}

	if len(fields) == 0 {
		logctx.Warn(ctx, "user not found by api key")
		return nil, models.ErrNotFound
	}

	userId, err := uuid.Parse(fields["id"])
	if err != nil {
		logctx.Error(ctx, "unexpected error parsing user id", logger.Error(err), logger.String("userId", fields["id"]))
		return nil, fmt.Errorf("unexpected error parsing user id: %w", err)
	}

	userType, err := models.StrToUserType(fields["type"])
	if err != nil {
		logctx.Error(ctx, "unexpected error parsing user type", logger.Error(err), logger.String("userId", userId.String()), logger.String("type", fields["type"]))
		return nil, fmt.Errorf("unexpected error parsing user type: %w", err)
	}

	if fields["apiKey"] != apiKey {
		logctx.Error(ctx, "api key mismatch", logger.String("userId", userId.String()))
		return nil, fmt.Errorf("api key mismatch")
	}

	logctx.Debug(ctx, "user found", logger.String("userId", userId.String()))

	return &models.User{
		Id:     userId,
		PubKey: fields["pubKey"],
		Type:   userType,
		ApiKey: fields["apiKey"],
	}, nil
}
