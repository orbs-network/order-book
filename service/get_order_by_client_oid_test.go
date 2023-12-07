package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_GetOrderByClientOId(t *testing.T) {

	ctx := context.Background()

	clientOId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")

	mockBcClient := &mocks.MockBcClient{IsVerified: true}

	t.Run("successfully retrieve order by client order ID - should return order", func(t *testing.T) {
		o := &models.Order{ClientOId: clientOId}
		svc, _ := service.New(&mocks.MockOrderBookStore{Order: o}, mockBcClient)

		order, err := svc.GetOrderByClientOId(ctx, clientOId)

		assert.NoError(t, err)
		assert.Equal(t, o, order)
	})

	t.Run("order not found - should return nil", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{Error: models.ErrNotFound}, mockBcClient)

		order, err := svc.GetOrderByClientOId(ctx, clientOId)

		assert.NoError(t, err)
		assert.Nil(t, order)
	})

	t.Run("unexpected error - should return error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{Error: assert.AnError}, mockBcClient)

		order, err := svc.GetOrderByClientOId(ctx, clientOId)

		assert.ErrorIs(t, err, assert.AnError)
		assert.Nil(t, order)
	})

}
