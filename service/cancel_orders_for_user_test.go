package service_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_CancelOrdersForUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully cancel all orders for a user", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: &mocks.User}

		s, _ := service.New(store)

		err := s.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Equal(t, err, nil)
	})

	t.Run("should return error on unexpected error", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: nil, Error: assert.AnError}

		s, _ := service.New(store)

		err := s.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "could not cancel orders for user")
	})

}
