package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_CancelOrder(t *testing.T) {

	userId := uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37")
	orderId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")
	order := &models.Order{UserId: userId}

	t.Run("no user in context - returns `ErrNoUserInContext` error", func(t *testing.T) {
		svc, _ := service.New(&mocks.MockOrderBookStore{})

		err := svc.CancelOrder(context.Background(), orderId)
		assert.Equal(t, models.ErrNoUserInContext, err)

	})

	testCases := []struct {
		name        string
		order       *models.Order
		err         error
		expectedErr error
	}{

		{name: "unexpected error when finding order - returns error", err: models.ErrOrderNotFound, expectedErr: models.ErrOrderNotFound},
		{name: "order not found - returns `ErrOrderNotFound` error", order: nil, err: nil, expectedErr: models.ErrOrderNotFound},
		{name: "user trying to cancel another user's order - returns `ErrUnauthorized` error", order: &models.Order{UserId: uuid.MustParse("00000000-0000-0000-0000-000000000009")}, expectedErr: models.ErrUnauthorized},
		{name: "unexpected error when removing order - returns error", order: order, err: assert.AnError, expectedErr: assert.AnError},
		{name: "order removed successfully - returns nil", order: order, err: nil, expectedErr: nil},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			fmt.Print(c.name)
			svc, _ := service.New(&mocks.MockOrderBookStore{Order: c.order, Error: c.err})

			userCtx := mocks.AddUserToCtx(&models.User{ID: userId, Type: models.MARKET_MAKER})

			err := svc.CancelOrder(userCtx, orderId)
			assert.Equal(t, c.expectedErr, err)
		})
	}

}
