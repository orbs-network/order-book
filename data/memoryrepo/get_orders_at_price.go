package memoryrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	priceStr := price.StringFixed(models.STR_PRECISION)
	ordersAtPrice, exists := r.sellOrders[symbol][priceStr]

	if !exists {
		return nil, fmt.Errorf("no orders found for symbol %s at price %s", symbol, priceStr)
	}

	orders := make([]models.Order, ordersAtPrice.List.Len())
	i := 0

	for e := ordersAtPrice.List.Front(); e != nil; e = e.Next() {
		orders[i] = e.Value.(models.Order)
		i++
	}

	return orders, nil
}
