package memoryrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
)

func (r *inMemoryRepository) RemoveOrder(ctx context.Context, orderId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderIdStr := orderId.String()

	// Retrieve the order's element from the orderLocations map
	element, exists := r.orderLocations[orderIdStr]
	if !exists {
		return fmt.Errorf("order with ID %s not found", orderIdStr)
	}

	order, ok := element.Value.(models.Order)
	if !ok {
		return fmt.Errorf("failed to cast element to Order type")
	}

	priceStr := order.Price.StringFixed(models.STR_PRECISION)

	// Remove the order from the orders at this price for this symbol
	ordersAtPrice := r.sellOrders[order.Symbol][priceStr]
	if ordersAtPrice == nil {
		return fmt.Errorf("no orders found at price %s for symbol %s", priceStr, order.Symbol)
	}

	ordersAtPrice.List.Remove(element)
	ordersAtPrice.Sum = ordersAtPrice.Sum.Sub(order.Size)

	// If no more orders at this price, remove the price entry
	if ordersAtPrice.List.Len() == 0 {
		delete(r.sellOrders[order.Symbol], priceStr)
	}

	// Remove the order from the user's orders
	if userOrders := r.userOrders[order.UserId.String()]; userOrders != nil {
		if symbolOrders := userOrders[order.Symbol]; symbolOrders != nil {
			delete(symbolOrders, priceStr)

			// If no more orders at this symbol, remove the symbol entry
			if len(symbolOrders) == 0 {
				delete(userOrders, order.Symbol)
			}
		}

		// If no more orders for this user, remove the user entry
		if len(userOrders) == 0 {
			delete(r.userOrders, order.UserId.String())
		}
	}

	// Remove the order from orderLocations map
	delete(r.orderLocations, orderIdStr)

	return nil
}
