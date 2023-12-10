package serviceuser

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServiceUser_GetUserByApiKey(t *testing.T) {
	ctx := context.Background()
	t.Run("should get user by api key", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{User: &mocks.User})

		user, err := userSvc.GetUserByApiKey(ctx, mocks.ApiKey)

		assert.NoError(t, err)
		assert.Equal(t, &mocks.User, user)

	})

	t.Run("should return error if failed to get user by api key", func(t *testing.T) {
		userSvc, _ := New(&mocks.MockUserStore{Error: assert.AnError})

		user, err := userSvc.GetUserByApiKey(ctx, mocks.ApiKey)

		assert.ErrorContains(t, err, "failed to get user by api key")
		assert.Nil(t, user)
	})
}
