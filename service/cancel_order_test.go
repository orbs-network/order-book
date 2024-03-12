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
	order := &models.Order{Id: orderId, UserId: userId, ClientOId: clientOId, Size: decimal.NewFromInt(9999999), SizePending: decimal.NewFromFloat(0), SizeFilled: decimal.NewFromFloat(50000)}

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
		{name: "order already cancelled - returns `ErrOrderCancelled` error", isClientOId: false, order: &models.Order{Cancelled: true}, expectedOrderId: nil, expectedErr: models.ErrOrderCancelled},
		{name: "order already filled so cannot be cancelled - returns `ErrOrderFilled`", isClientOId: false, order: &models.Order{UserId: userId, SizeFilled: decimal.NewFromFloat(99999.99), Size: decimal.NewFromFloat(99999.99)}, expectedOrderId: nil, expectedErr: models.ErrOrderFilled},
		{name: "order is partially filled and not pending", isClientOId: false, order: &models.Order{Id: order.Id, UserId: userId, SizeFilled: decimal.NewFromFloat(10), Size: decimal.NewFromFloat(50), SizePending: decimal.NewFromFloat(0)}, expectedOrderId: &order.Id, expectedErr: nil},
		{name: "order is partially filled and pending", isClientOId: false, order: &models.Order{Id: order.Id, UserId: userId, SizeFilled: decimal.NewFromFloat(10), Size: decimal.NewFromFloat(50), SizePending: decimal.NewFromFloat(40)}, expectedOrderId: &order.Id, expectedErr: nil},
		{name: "order is not filled and not pending", isClientOId: false, order: &models.Order{Id: order.Id, UserId: userId, SizeFilled: decimal.NewFromFloat(0), Size: decimal.NewFromFloat(50), SizePending: decimal.NewFromFloat(0)}, expectedOrderId: &order.Id, expectedErr: nil},
		{name: "order is not filled and pending", isClientOId: false, order: &models.Order{Id: order.Id, UserId: userId, SizeFilled: decimal.NewFromFloat(0), Size: decimal.NewFromFloat(50), SizePending: decimal.NewFromFloat(40)}, expectedOrderId: &order.Id, expectedErr: nil},
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
