package service_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_publishOrderEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("should subscribe to user order event updates", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{
			EventsChan: make(chan []byte),
		}, &mocks.MockBcClient{})

		channel, err := svc.SubscribeUserOrders(ctx, mocks.UserId)

		assert.NotNil(t, channel)
		assert.NoError(t, err)
	})

	t.Run("should return error when failed to subscribe to user order event updates", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{
			Error: assert.AnError,
		}, &mocks.MockBcClient{})

		channel, err := svc.SubscribeUserOrders(ctx, mocks.UserId)

		assert.Nil(t, channel)
		assert.Error(t, err)
	})
}
