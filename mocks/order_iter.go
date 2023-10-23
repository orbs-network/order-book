package mocks

import "github.com/orbs-network/order-book/models"

type OrderIter struct {
	Error          error
	Order          *models.Order
	ShouldHaveNext bool
}

func (o *OrderIter) HasNext() bool {
	return o.ShouldHaveNext
}

func (o *OrderIter) Next() *models.Order {
	return o.Order
}
