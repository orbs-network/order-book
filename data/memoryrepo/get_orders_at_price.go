package memoryrepo

import (
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) GetOrdersAtPrice(price decimal.Decimal) []models.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders, exists := r.sellOrders[price.StringFixed(models.STR_PRECISION)]
	if !exists {
		return nil
	}

	ordersSlice := make([]models.Order, orders.List.Len())
	i := 0
	for e := orders.List.Front(); e != nil; e = e.Next() {
		ordersSlice[i] = *e.Value.(*models.Order)
		i++
	}
	return ordersSlice
}
