package memoryrepo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
)

func (r *inMemoryRepository) RemoveOrder(orderId uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderIdStr := orderId.String()

	element, exists := r.orderLocations[orderIdStr]
	if !exists {
		return fmt.Errorf("order not found for order ID %q", orderIdStr)
	}

	order, ok := element.Value.(models.Order)
	if !ok {
		return fmt.Errorf("failed to cast order %q", orderIdStr)
	}

	orders := r.sellOrders[order.Price]

	orders.List.Remove(element)
	orders.Sum = orders.Sum.Sub(order.Size)

	// Remove the order from the orderLocations map
	delete(r.orderLocations, order.Id.String())

	// Remove the order from the userOrders map
	delete(r.userOrders[order.UserId.String()], order.Price)

	return nil
}

// func (r *inMemoryRepository) RemoveAllOrdersForUser(userId uuid.UUID) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	userIdStr := userId.String()

// 	// Get all orders for the user
// 	orders := r.userOrders[userIdStr]

// 	// Create a wait group to wait for all goroutines to complete
// 	var wg sync.WaitGroup

// 	// Remove all orders for the user
// 	for price, element := range orders {
// 		wg.Add(1)
// 		go func(price decimal.Decimal, element *list.Element) {
// 			defer wg.Done()

// 			order, ok := element.Value.(models.Order)
// 			if !ok {
// 				log.Printf("failed to cast order %q", order.Id.String())
// 				return
// 			}

// 			orders := r.sellOrders[price]

// 			orders.List.Remove(element)
// 			orders.Sum = orders.Sum.Sub(order.Size)

// 			// Remove the order from the orderLocations map
// 			delete(r.orderLocations, order.Id.String())
// 		}(price, element)
// 	}

// 	// Wait for all goroutines to complete
// 	wg.Wait()

// 	// Remove the user from the userOrders map
// 	delete(r.userOrders, userIdStr)

// 	return nil
// }
