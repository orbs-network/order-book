package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type OrderIterMock struct {
	Orders []models.Order
	Index  int
}

func (i *OrderIterMock) HasNext() bool {
	return i.Index < (len(i.Orders) - 1)
}

func (i *OrderIterMock) Next(ctx context.Context) *models.Order {
	if i.Index >= len(i.Orders) {
		logctx.Error(ctx, "Error iterator reached last element")
		return nil
	}

	i.Index = i.Index + 1

	return &i.Orders[i.Index]
}
