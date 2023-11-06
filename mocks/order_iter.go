package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type OrderIterMock struct {
	orders []models.Order
	index  int
}

func (i *OrderIterMock) HasNext() bool {
	return i.index < (len(i.orders) - 1)
}

func (i *OrderIterMock) Next(ctx context.Context) *models.Order {
	if i.index >= len(i.orders) {
		logctx.Error(ctx, "Error iterator reached last element")
		return nil
	}

	// increment index
	i.index = i.index + 1

	return &i.orders[i.index]
}
