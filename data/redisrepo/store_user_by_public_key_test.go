package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreUserByPublicKey(t *testing.T) {

	t.Run("should successfully store user details", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPubKeyKey(mocks.PubKey)

		mock.ExpectHMSet(key, map[string]interface{}{
			"id":     mocks.UserId.String(),
			"pubKey": mocks.PubKey,
			"type":   mocks.UserType.String(),
		}).SetVal(true)

		err := repo.StoreUserByPublicKey(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
		})

		assert.NoError(t, err)
	})

	t.Run("should return error on unexpected error storing user by public key", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateUserPubKeyKey(mocks.PubKey)

		mock.ExpectHMSet(key, map[string]interface{}{
			"id":     mocks.UserId.String(),
			"pubKey": mocks.PubKey,
			"type":   mocks.UserType.String(),
		}).SetErr(assert.AnError)

		err := repo.StoreUserByPublicKey(ctx, models.User{
			Id:     mocks.UserId,
			PubKey: mocks.PubKey,
			Type:   mocks.UserType,
		})

		assert.ErrorContains(t, err, "unexpected error storing user by public key")
	})

}
