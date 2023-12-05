package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_CreateUser(t *testing.T) {

	fields := map[string]interface{}{
		"id":     mocks.UserId.String(),
		"type":   mocks.UserType.String(),
		"pubKey": mocks.PubKey,
		"apiKey": mocks.ApiKey,
	}

	t.Run("should successfully create user and return user instance on success", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		userApiKey := CreateUserApiKeyKey(mocks.ApiKey)
		userIdKey := CreateUserIdKey(mocks.UserId)

		mock.ExpectExists(userApiKey, userIdKey).SetVal(0)
		mock.ExpectTxPipeline()
		mock.ExpectHMSet(userApiKey, fields).SetVal(true)
		mock.ExpectHMSet(userIdKey, fields).SetVal(true)
		mock.ExpectTxPipelineExec()

		user, err := repo.CreateUser(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		})

		assert.Equal(t, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		}, user)
		assert.NoError(t, err)

	})

	t.Run("should return `ErrUserAlreadyExists` error if user already exists", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		userApiKey := CreateUserApiKeyKey(mocks.ApiKey)
		userIdKey := CreateUserIdKey(mocks.UserId)

		mock.ExpectExists(userApiKey, userIdKey).SetVal(1)

		user, err := repo.CreateUser(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		})

		assert.Equal(t, models.User{}, user)
		assert.ErrorIs(t, err, models.ErrUserAlreadyExists)
	})

	t.Run("should return error on unexpected exists error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		userApiKey := CreateUserApiKeyKey(mocks.ApiKey)
		userIdKey := CreateUserIdKey(mocks.UserId)

		mock.ExpectExists(userApiKey, userIdKey).SetErr(assert.AnError)

		user, err := repo.CreateUser(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		})

		assert.Equal(t, models.User{}, user)
		assert.ErrorContains(t, err, "unexpected error checking if user exists")
	})

	t.Run("should return error on unexpected create user error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		userApiKey := CreateUserApiKeyKey(mocks.ApiKey)
		userIdKey := CreateUserIdKey(mocks.UserId)

		mock.ExpectExists(userApiKey, userIdKey).SetVal(0)
		mock.ExpectTxPipeline()
		mock.ExpectHMSet(userApiKey, fields).SetErr(assert.AnError)
		mock.ExpectTxPipelineExec()

		user, err := repo.CreateUser(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		})

		assert.Equal(t, models.User{}, user)
		assert.ErrorContains(t, err, "unexpected error creating user")
	})

}
