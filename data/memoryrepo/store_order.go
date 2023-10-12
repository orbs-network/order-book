package memoryrepo

import (
	"container/list"
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) StoreOrder(ctx context.Context, order models.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userIdStr := order.UserId.String()
	priceStr := order.Price.StringFixed(models.STR_PRECISION)
	symbol := order.Symbol

	// Initilize if no previous orders for this user
	if r.userOrders[userIdStr] == nil {
		r.userOrders[userIdStr] = make(map[models.Symbol]map[string]*list.Element)
	}

	// Initilize if no previous orders for this symbol for this user
	if r.userOrders[userIdStr][symbol] == nil {
		r.userOrders[userIdStr][symbol] = make(map[string]*list.Element)
	}

	if r.userOrders[userIdStr][symbol][priceStr] != nil {
		return fmt.Errorf("user %s already has an order at price %s for symbol %s", userIdStr, priceStr, symbol)
	}

	// Initialize if no previous orders for this symbol
	if r.sellOrders[symbol] == nil {
		r.sellOrders[symbol] = make(map[string]*ordersAtPrice)
	}

	orders, exists := r.sellOrders[symbol][priceStr]
	if !exists {
		orders = &ordersAtPrice{List: list.New(), Sum: decimal.Zero}
		r.sellOrders[symbol][priceStr] = orders
	}

	element := orders.List.PushBack(order)
	orders.Sum = orders.Sum.Add(order.Size)

	r.orderLocations[order.Id.String()] = element

	r.userOrders[userIdStr][symbol][priceStr] = element

	return nil
}
