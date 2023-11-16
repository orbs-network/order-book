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
	clientOId := uuid.MustParse("f577273e-12de-4acc-a4f8-de7fb5b86e37")
	order := &models.Order{Id: orderId, UserId: userId, Status: models.STATUS_OPEN, ClientOId: clientOId}

	testCases := []struct {
		name            string
		order           *models.Order
		isClientOId     bool
		err             error
		expectedOrderId *uuid.UUID
		expectedErr     error
	}{
		{name: "unexpected error when finding order by orderId - returns error", isClientOId: false, err: assert.AnError, expectedOrderId: nil, expectedErr: assert.AnError},
		{name: "unexpected error when finding order by clientOId - returns error", isClientOId: true, err: assert.AnError, expectedOrderId: nil, expectedErr: assert.AnError},
		{name: "order not found - returns `ErrOrderNotFound` error", isClientOId: false, order: nil, err: nil, expectedOrderId: nil, expectedErr: models.ErrOrderNotFound},
		// {name: "trying to cancel order that is not open - returns `ErrOrderNotOpen` error", isClientOId: false, order: &models.Order{Status: models.STATUS_PENDING}, err: nil, expectedOrderId: nil, expectedErr: models.ErrOrderNotOpen},
		{name: "user trying to cancel another user's order - returns `ErrUnauthorized` error", isClientOId: false, order: &models.Order{UserId: uuid.MustParse("00000000-0000-0000-0000-000000000009"), Status: models.STATUS_OPEN}, expectedOrderId: nil, expectedErr: models.ErrUnauthorized},
		{name: "unexpected error when removing order - returns error", isClientOId: false, order: order, err: assert.AnError, expectedOrderId: nil, expectedErr: assert.AnError},
		{name: "order removed successfully by orderId - returns cancelled orderId", isClientOId: false, order: order, err: nil, expectedOrderId: &orderId, expectedErr: nil},
		{name: "order removed successfully by clientOId - returns cancelled orderId", isClientOId: true, order: order, err: nil, expectedOrderId: &orderId, expectedErr: nil},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			fmt.Print(c.name)
			svc, _ := service.New(&mocks.MockOrderBookStore{Order: c.order, Error: c.err})

			input := service.CancelOrderInput{Id: orderId, IsClientOId: false, UserId: userId}

			orderId, err := svc.CancelOrder(context.Background(), input)
			assert.Equal(t, c.expectedOrderId, orderId)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
