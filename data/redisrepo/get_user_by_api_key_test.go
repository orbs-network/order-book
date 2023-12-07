package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetUserByPublicKey(t *testing.T) {

	mockApiKey := "mock-api-key"

	t.Run("should return user by public key", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"type":   mocks.UserType.String(),
			"pubKey": mocks.PubKey,
			"apiKey": mockApiKey,
		})

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.NoError(t, err)
		assert.Equal(t, user, &models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mockApiKey,
		})
	})

	t.Run("should return user not found error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetVal(map[string]string{})

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("should return error on unexpected error getting user by api key", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetErr(assert.AnError)

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error getting user by api key")
	})

	t.Run("should return error on unexpected error parsing user id", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":     "invalid",
			"pubKey": mocks.PubKey,
			"type":   mocks.UserType.String(),
			"apiKey": mockApiKey,
		})

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error parsing user id")
	})

	t.Run("should return error on unexpected error parsing user type", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"pubKey": mocks.PubKey,
			"type":   "invalid",
			"apiKey": mockApiKey,
		})

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error parsing user type")
	})

	t.Run("should return error on api key mismatch", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserApiKeyKey(mockApiKey)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"pubKey": mocks.PubKey,
			"type":   mocks.UserType.String(),
			"apiKey": "a-different-api-key",
		})

		user, err := repo.GetUserByApiKey(ctx, mockApiKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "api key mismatch")
	})

}
