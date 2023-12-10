package serviceuser

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceUser_GetUserById(t *testing.T) {
	ctx := context.Background()

	t.Run("should get user by id", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		user, err := userSvc.GetUserById(ctx, mocks.UserId)

		assert.NoError(t, err)
		assert.Equal(t, &mocks.User, user)
	})

	t.Run("should return `ErrNotFound` error if user not found", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{Error: models.ErrNotFound})

		user, err := userSvc.GetUserById(ctx, mocks.UserId)

		assert.ErrorIs(t, err, models.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("should return error if failed to get user by id", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{Error: assert.AnError})

		user, err := userSvc.GetUserById(ctx, mocks.UserId)

		assert.ErrorContains(t, err, "failed to get user by id")
		assert.Nil(t, user)
	})
}
