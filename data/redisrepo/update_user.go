package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// UpdateUser updates a user's pubKey and apiKey
//
// It deletes the user by id and apiKey and creates a new user by id and (new) apiKey (with the updated pubKey)
//
// Ensure that `PubKey` and `ApiKey` are correct and not empty
func (r *redisRepository) UpdateUser(ctx context.Context, input storeuser.UpdateUserInput) error {

	if input.ApiKey == "" || input.PubKey == "" {
		logctx.Error(ctx, "apiKey or pubKey is empty", logger.String("apiKey", input.ApiKey), logger.String("pubKey", input.PubKey))
		return models.ErrInvalidInput
	}

	user, err := r.GetUserById(ctx, input.UserId)
	if err != nil {
		logctx.Error(ctx, "unexpected error getting user by id", logger.Error(err), logger.String("userId", input.UserId.String()))
		return fmt.Errorf("unexpected error getting user by id: %w", err)
	}

	updatedFields := map[string]interface{}{
		"id":     user.Id.String(),
		"type":   user.Type.String(),
		"pubKey": input.PubKey,
		"apiKey": input.ApiKey,
	}

	rUserIdKey := CreateUserIdKey(user.Id)
	rOldUserApiKey := CreateUserApiKeyKey(user.ApiKey)
	rNewUserApiKey := CreateUserApiKeyKey(input.ApiKey)

	// --- START TRANSACTION ---
	transaction := r.client.TxPipeline()
	// delete user by id
	transaction.Del(ctx, rUserIdKey)
	// delete user by api key
	transaction.Del(ctx, rOldUserApiKey)
	// create user by id
	transaction.HMSet(ctx, rUserIdKey, updatedFields)
	// create user by api key
	transaction.HMSet(ctx, rNewUserApiKey, updatedFields)
	_, err = transaction.Exec(ctx)
	// --- END TRANSACTION ---

	if err != nil {
		logctx.Error(ctx, "transaction failed updating user", logger.Error(err), logger.String("userId", input.UserId.String()), logger.String("pubKey", input.PubKey))
		return fmt.Errorf("transaction failed: %w", err)
	}

	logctx.Info(ctx, "user updated", logger.String("userId", input.UserId.String()), logger.String("pubKey", input.PubKey))
	return nil
}
