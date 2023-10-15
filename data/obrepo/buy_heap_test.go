package obrepo

import (
	"container/heap"
	"container/list"
	"fmt"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Test that the order with the highest price is at the root of the heap.
func TestBuyHeapOrdering(t *testing.T) {
	h := buyHeap{}
	heap.Init(&h)

	order_one := &ordersAtPrice{Price: decimal.NewFromFloat(12.0), Orders: list.New()}
	order_one.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(12.0)})

	order_two := &ordersAtPrice{Price: decimal.NewFromFloat(99.0), Orders: list.New()}
	order_two.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(99.0)})

	order_three := &ordersAtPrice{Price: decimal.NewFromFloat(1.0), Orders: list.New()}
	order_three.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(1.0)})

	heap.Push(&h, order_one)
	heap.Push(&h, order_two)
	heap.Push(&h, order_three)

	max := heap.Pop(&h).(*ordersAtPrice)
	assert.Equal(t, decimal.NewFromFloat(99.0), max.Price, "Price should be 99.0 as buy orders are in descending order")
}

func TestBuyHeapRemovalByIndex(t *testing.T) {
	h := buyHeap{}
	heap.Init(&h)

	order_one := &ordersAtPrice{Price: decimal.NewFromFloat(12.0), Orders: list.New()}
	order_one.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(12.0)})

	order_two := &ordersAtPrice{Price: decimal.NewFromFloat(99.0), Orders: list.New()}
	order_two.Orders.PushBack(models.Order{Price: decimal.NewFromFloat(99.0)})

	heap.Push(&h, order_one)
	heap.Push(&h, order_two)
	fmt.Printf("order_one.Index: %#v\n", order_one.Index)

	heap.Remove(&h, order_one.Index)

	assert.Equal(t, 1, h.Len(), "Heap should have 1 item after removal")
	assert.Equal(t, decimal.NewFromFloat(99.0), h[0].Price, "Price should be 99.0 as buy orders are in descending order")
}
