package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) GetUserByPublicKey(ctx context.Context, publicKey string) (*models.User, error) {

	key := CreateUserPubKeyKey(publicKey)

	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		logctx.Error(ctx, "unexpected error getting user by public key", logger.Error(err), logger.String("publicKey", publicKey))
		return nil, fmt.Errorf("unexpected error getting user by public key: %w", err)
	}

	if len(fields) == 0 {
		logctx.Warn(ctx, "user not found", logger.String("publicKey", publicKey))
		return nil, models.ErrUserNotFound
	}

	userId, err := uuid.Parse(fields["id"])
	if err != nil {
		logctx.Error(ctx, "unexpected error parsing user id", logger.Error(err), logger.String("publicKey", publicKey), logger.String("userId", fields["id"]))
		return nil, fmt.Errorf("unexpected error parsing user id: %w", err)
	}

	userType, err := models.StrToUserType(fields["type"])
	if err != nil {
		logctx.Error(ctx, "unexpected error parsing user type", logger.Error(err), logger.String("publicKey", publicKey), logger.String("type", fields["type"]))
		return nil, fmt.Errorf("unexpected error parsing user type: %w", err)
	}

	if fields["pubKey"] != publicKey {
		logctx.Error(ctx, "public key mismatch", logger.String("publicKey from args", publicKey), logger.String("pubKey from map", fields["pubKey"]))
		return nil, fmt.Errorf("public key mismatch")
	}

	logctx.Info(ctx, "user found", logger.String("publicKey", publicKey), logger.String("userId", userId.String()))

	return &models.User{
		Id:     userId,
		PubKey: publicKey,
		Type:   userType,
	}, nil
}
