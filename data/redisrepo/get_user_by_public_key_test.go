package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetUserByPublicKey(t *testing.T) {

	t.Run("should return user by public key", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":   mocks.UserId.String(),
			"pk":   mocks.Pk,
			"type": mocks.UserType.String(),
		})

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.NoError(t, err)
		assert.Equal(t, user, &models.User{
			Id:   mocks.UserId,
			Pk:   mocks.Pk,
			Type: mocks.UserType,
		})
	})

	t.Run("should return user not found error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetVal(map[string]string{})

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, models.ErrUserNotFound)
	})

	t.Run("should return error on unexpected error getting user by public key", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetErr(assert.AnError)

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error getting user by public key")
	})

	t.Run("should return error on unexpected error parsing user id", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":   "invalid",
			"pk":   mocks.Pk,
			"type": mocks.UserType.String(),
		})

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error parsing user id")
	})

	t.Run("should return error on unexpected error parsing user type", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":   mocks.UserId.String(),
			"pk":   mocks.Pk,
			"type": "invalid",
		})

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error parsing user type")
	})

	t.Run("should return error on public key mismatch", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPKKey(mocks.Pk)

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"id":   mocks.UserId.String(),
			"pk":   "invalid",
			"type": mocks.UserType.String(),
		})

		user, err := repo.GetUserByPublicKey(ctx, mocks.Pk)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "public key mismatch")
	})

}
