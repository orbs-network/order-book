package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreUserByPublicKey(ctx context.Context, user models.User) error {

	key := CreateUserPubKeyKey(user.PubKey)

	fields := map[string]interface{}{
		"id":     user.Id.String(),
		"type":   user.Type.String(),
		"pubKey": user.PubKey,
	}

	_, err := r.client.HMSet(ctx, key, fields).Result()
	if err != nil {
		logctx.Error(ctx, "unexpected error storing user by public key", logger.Error(err), logger.String("publicKey", user.PubKey))
		return fmt.Errorf("unexpected error storing user by public key: %w", err)
	}

	logctx.Info(ctx, "user stored by public key", logger.String("publicKey", user.PubKey), logger.String("userId", user.Id.String()), logger.String("type", user.Type.String()))

	return nil
}
