package memoryrepo

import (
	"container/list"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) StoreOrder(order models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIdStr := order.UserId.String()
	priceStr := order.Price.StringFixed(models.STR_PRECISION)

	orders, exists := r.sellOrders[priceStr]
	if !exists {
		orders = &ordersAtPrice{List: list.New(), Sum: decimal.Zero}
		r.sellOrders[priceStr] = orders
	}

	element := orders.List.PushBack(order)
	orders.Sum = orders.Sum.Add(order.Size)

	r.orderLocations[order.Id.String()] = element

	// Record the order under the user's ID and price
	if r.userOrders[userIdStr] == nil {
		r.userOrders[userIdStr] = make(map[string]*list.Element)
	}

	r.userOrders[userIdStr][priceStr] = element

	return nil
}
