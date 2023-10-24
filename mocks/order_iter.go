package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

type OrderIter struct {
	Error          error
	Order          *models.Order
	ShouldHaveNext bool
}

func (o *OrderIter) HasNext() bool {
	return o.ShouldHaveNext
}

func (o *OrderIter) Next(ctx context.Context) *models.Order {
	return o.Order
}
