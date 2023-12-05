package serviceuser

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestServiceUser_UpdateUser(t *testing.T) {
	ctx := context.Background()
	t.Run("should update user", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		err := userSvc.UpdateUser(ctx, UpdateUserInput{
			UserId: mocks.UserId,
			ApiKey: mocks.ApiKey,
			PubKey: mocks.PubKey,
		})

		assert.NoError(t, err)
	})

	t.Run("should return error if no api key", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		err := userSvc.UpdateUser(ctx, UpdateUserInput{
			UserId: mocks.UserId,
			ApiKey: "",
			PubKey: mocks.PubKey,
		})

		assert.ErrorIs(t, err, models.ErrInvalidInput)
	})

	t.Run("should return error if no pub key", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		err := userSvc.UpdateUser(ctx, UpdateUserInput{
			UserId: mocks.UserId,
			ApiKey: mocks.ApiKey,
			PubKey: "",
		})

		assert.ErrorIs(t, err, models.ErrInvalidInput)
	})

	t.Run("should return error if failed to update user", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{Error: assert.AnError})

		err := userSvc.UpdateUser(ctx, UpdateUserInput{
			UserId: mocks.UserId,
			ApiKey: mocks.ApiKey,
			PubKey: mocks.PubKey,
		})

		assert.ErrorContains(t, err, "failed to update user")
	})
}
