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
	mockBcClient := &mocks.MockBcClient{IsVerified: true}

	t.Run("should successfully cancel all orders for a user", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: &mocks.User, Orders: []models.Order{mocks.Order}, Order: &mocks.Order}

		s, _ := service.New(store, mockBcClient)

		orderIds, err := s.CancelOrdersForUser(ctx, mocks.UserId, "")
		assert.Equal(t, orderIds[0], mocks.Order.Id)
		assert.Equal(t, err, nil)
	})

	t.Run("should return error if no orders found for user", func(t *testing.T) {
		store := &mocks.MockOrderBookStore{User: &mocks.User, Orders: []models.Order{}, Error: models.ErrNotFound}

		s, _ := service.New(store, mockBcClient)

		orderIds, err := s.CancelOrdersForUser(ctx, mocks.UserId, "")
		assert.Empty(t, orderIds)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	// t.Run("should return error on unexpected error", func(t *testing.T) {
	// 	store := &mocks.MockOrderBookStore{User: nil, Error: assert.AnError, Order: &mocks.Order}

	// 	s, _ := service.New(store, mockBcClient)

	// 	orderIds, err := s.CancelOrdersForUser(ctx, mocks.UserId, "")
	// 	assert.Empty(t, orderIds)
	// 	assert.ErrorContains(t, err, "could not cancel orders for user")
	// })

}
