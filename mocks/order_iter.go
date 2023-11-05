package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type OrderIterMock struct {
	orders []models.Order
	index  int
	// Error          error
	// Order          *models.Order
	// ShouldHaveNext bool
}

func (i OrderIterMock) HasNext() bool {
	return i.index < (len(i.orders) - 1)

}

func (i OrderIterMock) Next(ctx context.Context) *models.Order {
	if i.index >= len(i.orders) {
		logctx.Error(ctx, "Error iterator reached last element")
		return nil
	}

	// increment index
	i.index = i.index + 1
	// get order

	return &i.orders[i.index]
}
