package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetUserById(t *testing.T) {

	t.Run("should return user by ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserIdKey(mocks.UserId)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":     mocks.UserId.String(),
			"type":   mocks.UserType.String(),
			"pubKey": mocks.PubKey,
			"apiKey": mocks.ApiKey,
		})

		user, err := repo.GetUserById(ctx, mocks.UserId)

		assert.NoError(t, err)
		assert.Equal(t, user, &models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
			ApiKey: mocks.ApiKey,
		})
	})

	t.Run("should return user not found error when no user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserIdKey(mocks.UserId)

		mock.ExpectHGetAll(key).SetVal(map[string]string{})

		_, err := repo.GetUserById(ctx, mocks.UserId)

		assert.Error(t, err)
		assert.Equal(t, err, models.ErrNotFound)
	})

	t.Run("should return error on unexpected error getting user by ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserIdKey(mocks.UserId)

		mock.ExpectHGetAll(key).SetErr(assert.AnError)

		_, err := repo.GetUserById(ctx, mocks.UserId)

		assert.Error(t, err)
		assert.Equal(t, err, assert.AnError)
	})
}
