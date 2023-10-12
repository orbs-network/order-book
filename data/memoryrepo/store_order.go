package memoryrepo

import (
	"container/list"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) StoreOrder(order models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIdStr := order.UserId.String()

	// Check if the user has already placed an order at this price
	if r.userOrders[userIdStr][order.Price] != nil {
		return fmt.Errorf("user %s already has an order at price %s", userIdStr, order.Price)
	}

	orders, exists := r.sellOrders[order.Price]
	if !exists {
		orders = &ordersAtPrice{List: list.New(), Sum: decimal.Zero}
		r.sellOrders[order.Price] = orders
	}

	element := orders.List.PushBack(order)
	orders.Sum = orders.Sum.Add(order.Size)

	r.orderLocations[order.Id.String()] = element

	// Record the order under the user's ID and price
	if r.userOrders[userIdStr] == nil {
		r.userOrders[userIdStr] = make(map[decimal.Decimal]*list.Element)
	}
	r.userOrders[userIdStr][order.Price] = element

	return nil
}
