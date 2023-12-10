package serviceuser

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceUser_CreateUser(t *testing.T) {
	ctx := context.Background()
	t.Run("should create user and return user instance on success", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		user, err := userSvc.CreateUser(ctx, CreateUserInput{
			PubKey: mocks.PubKey,
		})

		assert.Equal(t, mocks.User, user)
		assert.NoError(t, err)

	})

	t.Run("should return `ErrUserAlreadyExists` error if user already exists", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &models.User{}, Error: models.ErrUserAlreadyExists})

		_, err := userSvc.CreateUser(ctx, CreateUserInput{
			PubKey: mocks.PubKey,
		})

		// assert.Equal(t, models.User{}, user)
		assert.ErrorIs(t, err, models.ErrUserAlreadyExists)
	})

	t.Run("should return error if failed to create user", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &models.User{}, Error: assert.AnError})

		user, err := userSvc.CreateUser(ctx, CreateUserInput{
			PubKey: mocks.PubKey,
		})

		assert.Equal(t, models.User{}, user)
		assert.ErrorContains(t, err, "failed to create user")
	})
}
