package service_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_CancelOrdersForUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully cancel all orders for a user", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: &mocks.User}

		s, _ := service.New(store)

		err := s.CancelOrdersForUser(ctx, mocks.PubKey)
		assert.Equal(t, err, nil)
	})

	t.Run("should return user not found error when no user found", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, Error: models.ErrUserNotFound}

		s, _ := service.New(store)

		err := s.CancelOrdersForUser(ctx, mocks.PubKey)
		assert.Equal(t, err, models.ErrUserNotFound)
	})

	t.Run("should return error on unexpected error getting user by public key", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, Error: assert.AnError}

		s, _ := service.New(store)

		err := s.CancelOrdersForUser(ctx, mocks.PubKey)
		assert.ErrorContains(t, err, "unexpected error getting user by public key")
	})

}
