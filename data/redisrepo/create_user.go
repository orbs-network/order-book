package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// CreateUser adds 2 entries to Redis:
// 1. User data indexed by API key
// 2. User data indexed by userId
func (r *redisRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {

	apiKeyKey := CreateUserApiKeyKey(user.ApiKey)
	userIdKey := CreateUserIdKey(user.Id)

	if exists, err := r.client.Exists(ctx, apiKeyKey, userIdKey).Result(); err != nil {
		logctx.Error(ctx, "unexpected error checking if user exists", logger.String("userId", user.Id.String()), logger.Error(err))
		return models.User{}, fmt.Errorf("unexpected error checking if user exists: %w", err)
	} else if exists > 0 {
		logctx.Warn(ctx, "user already exists", logger.String("userId", user.Id.String()))
		return models.User{}, models.ErrUserAlreadyExists
	}

	fields := map[string]interface{}{
		"id":     user.Id.String(),
		"type":   user.Type.String(),
		"pubKey": user.PubKey,
		"apiKey": user.ApiKey,
	}

	transaction := r.client.TxPipeline()
	transaction.HMSet(ctx, apiKeyKey, fields)
	transaction.HMSet(ctx, userIdKey, fields)
	_, err := transaction.Exec(ctx)

	if err != nil {
		logctx.Error(ctx, "unexpected error creating user", logger.String("userId", user.Id.String()))
		return models.User{}, fmt.Errorf("unexpected error creating user: %w", err)
	}

	logctx.Info(ctx, "user created", logger.String("userId", user.Id.String()))
	return user, nil
}
