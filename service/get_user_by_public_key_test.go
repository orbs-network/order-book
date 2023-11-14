package service_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUserByPublicKey(t *testing.T) {

	ctx := context.Background()

	t.Run("should get a user by their public key", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: &mocks.User}

		svc, _ := service.New(store)

		user, _ := svc.GetUserByPublicKey(ctx, mocks.PubKey)

		assert.Equal(t, user, &mocks.User)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, ErrUser: models.ErrUserNotFound}

		svc, _ := service.New(store)

		user, err := svc.GetUserByPublicKey(ctx, mocks.PubKey)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, models.ErrUserNotFound)
	})

	t.Run("should return error on unexpected error getting user by public key", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, ErrUser: assert.AnError}

		svc, _ := service.New(store)

		user, err := svc.GetUserByPublicKey(ctx, mocks.PubKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "unexpected error getting user by public key")
	})

	t.Run("should return error if user is nil but no error", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, ErrUser: nil}

		svc, _ := service.New(store)

		user, err := svc.GetUserByPublicKey(ctx, mocks.PubKey)

		assert.Nil(t, user)
		assert.ErrorContains(t, err, "user is nil but no error returned")
	})
}
