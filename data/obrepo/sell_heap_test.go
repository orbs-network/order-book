package obrepo

import (
	"container/heap"
	"container/list"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Test that the order with the lowest price is at the root of the heap.
func TestSellHeapOrdering(t *testing.T) {
	h := sellHeap{}

	order_one := &ordersAtPrice{Price: decimal.NewFromFloat(12.0), Orders: list.New()}
	order_one.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(12.0)})

	order_two := &ordersAtPrice{Price: decimal.NewFromFloat(99.0), Orders: list.New()}
	order_two.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(99.0)})

	order_three := &ordersAtPrice{Price: decimal.NewFromFloat(1.0), Orders: list.New()}
	order_three.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(1.0)})

	h.Push(order_one)
	h.Push(order_two)
	h.Push(order_three)

	heap.Init(&h)

	rootNode := *h[0]

	assert.Equal(t, decimal.NewFromFloat(1.0), rootNode.Price, "Price should be 1.0 as sell orders are in ascending order")
}
