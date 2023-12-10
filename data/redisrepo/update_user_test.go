package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_UpdateUser(t *testing.T) {
	oldApiKey := "old-api-key"
	newPubKey := "new-pub-key"
	newApiKey := "new-api-key"

	input := storeuser.UpdateUserInput{
		UserId: mocks.UserId,
		PubKey: newPubKey,
		ApiKey: newApiKey,
	}

	fields := map[string]interface{}{
		"id":     mocks.UserId.String(),
		"type":   mocks.UserType.String(),
		"pubKey": newPubKey,
		"apiKey": newApiKey,
	}

	t.Run("should update user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		rUserIdKey := CreateUserIdKey(mocks.UserId)
		rOldUserApiKey := CreateUserApiKeyKey(oldApiKey)
		rNewUserApiKey := CreateUserApiKeyKey(newApiKey)

		mock.ExpectHGetAll(rUserIdKey).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"type":   mocks.UserType.String(),
			"pubKey": mocks.PubKey,
			"apiKey": oldApiKey,
		})

		mock.ExpectTxPipeline()
		mock.ExpectDel(rUserIdKey).SetVal(1)
		mock.ExpectDel(rOldUserApiKey).SetVal(1)
		mock.ExpectHMSet(rUserIdKey, fields).SetVal(true)
		mock.ExpectHMSet(rNewUserApiKey, fields).SetVal(true)
		mock.ExpectTxPipelineExec()

		err := repo.UpdateUser(ctx, input)

		assert.NoError(t, err)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		rUserIdKey := CreateUserIdKey(mocks.UserId)

		mock.ExpectHGetAll(rUserIdKey).SetVal(map[string]string{})

		err := repo.UpdateUser(ctx, input)

		assert.Error(t, err)
	})

	t.Run("should return error if transaction failed", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		rUserIdKey := CreateUserIdKey(mocks.UserId)
		rOldUserApiKey := CreateUserApiKeyKey(oldApiKey)
		rNewUserApiKey := CreateUserApiKeyKey(newApiKey)

		mock.ExpectHGetAll(rUserIdKey).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"type":   mocks.UserType.String(),
			"pubKey": mocks.PubKey,
			"apiKey": oldApiKey,
		})

		mock.ExpectTxPipeline()
		mock.ExpectDel(rUserIdKey).SetVal(1)
		mock.ExpectDel(rOldUserApiKey).SetVal(1)
		mock.ExpectHMSet(rUserIdKey, fields).SetVal(true)
		mock.ExpectHMSet(rNewUserApiKey, fields).SetErr(assert.AnError)
		mock.ExpectTxPipelineExec()

		err := repo.UpdateUser(ctx, input)

		assert.ErrorContains(t, err, "transaction failed")
	})
}
