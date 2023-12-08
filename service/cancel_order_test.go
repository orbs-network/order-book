package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestService_CancelOrder(t *testing.T) {

	userId := uuid.MustParse("a577273e-12de-4acc-a4f8-de7fb5b86e37")
	orderId := uuid.MustParse("e577273e-12de-4acc-a4f8-de7fb5b86e37")
	clientOId := uuid.MustParse("f577273e-12de-4acc-a4f8-de7fb5b86e37")
	order := &models.Order{Id: orderId, UserId: userId, ClientOId: clientOId}

	mockBcClient := &mocks.MockBcClient{IsVerified: true}

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
		{name: "order not found - returns `ErrNotFound` error", isClientOId: false, order: nil, err: nil, expectedOrderId: nil, expectedErr: models.ErrNotFound},
		{name: "cancelling order not possible when order is pending", isClientOId: false, order: &models.Order{UserId: userId, SizePending: decimal.NewFromFloat(254), SizeFilled: decimal.NewFromFloat(32323.32)}, expectedOrderId: nil, expectedErr: models.ErrOrderPending},
		{name: "unexpected error when removing order - returns error", isClientOId: false, order: order, err: assert.AnError, expectedOrderId: nil, expectedErr: assert.AnError},
		{name: "order removed successfully by orderId - returns cancelled orderId", isClientOId: false, order: order, err: nil, expectedOrderId: &orderId, expectedErr: nil},
		{name: "order removed successfully by clientOId - returns cancelled orderId", isClientOId: true, order: order, err: nil, expectedOrderId: &orderId, expectedErr: nil},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			fmt.Print(c.name)
			svc, _ := service.New(&mocks.MockOrderBookStore{Order: c.order, Error: c.err}, mockBcClient)

			input := service.CancelOrderInput{Id: orderId, IsClientOId: false, UserId: userId}

			orderId, err := svc.CancelOrder(context.Background(), input)
			assert.Equal(t, c.expectedOrderId, orderId)
			assert.Equal(t, c.expectedErr, err)
		})
	}
}
