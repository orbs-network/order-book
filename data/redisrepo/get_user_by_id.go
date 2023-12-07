package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// GetUserById returns a user by their userId
func (r *redisRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {

	key := CreateUserIdKey(userId)

	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		logctx.Error(ctx, "user not found by ID", logger.String("userId", userId.String()))
		return nil, models.ErrUserNotFound
	}

	userType, err := models.StrToUserType(fields["type"])
	if err != nil {
		logctx.Error(ctx, "unexpected error parsing user type", logger.Error(err), logger.String("userId", userId.String()), logger.String("type", fields["type"]))
		return nil, fmt.Errorf("unexpected error parsing user type: %w", err)
	}

	return &models.User{
		Id:     userId,
		PubKey: fields["pubKey"],
		Type:   userType,
		ApiKey: fields["apiKey"],
	}, nil
}
