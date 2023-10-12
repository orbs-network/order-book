package memoryrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// FindOrder returns the order (if any) for the given user, symbol and price
func (r *inMemoryRepository) FindOrder(ctx context.Context, input models.FindOrderInput) (*models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	priceStr := input.Price.StringFixed(models.STR_PRECISION)

	// Retrieve the user's orders
	userOrders, userExists := r.userOrders[input.UserId.String()]
	if !userExists {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", input.UserId.String()))
		return nil, nil
	}

	// Retrieve the order for the specified symbol
	symbolOrders, symbolExists := userOrders[input.Symbol]
	if !symbolExists {
		logctx.Info(ctx, "no orders found for symbol", logger.String("userId", input.UserId.String()), logger.String("symbol", input.Symbol.String()))
		return nil, nil
	}

	// Retrieve the order element for the specified price
	element, priceExists := symbolOrders[priceStr]
	if !priceExists {
		logctx.Info(ctx, "no orders found for price", logger.String("userId", input.UserId.String()), logger.String("symbol", input.Symbol.String()), logger.String("price", priceStr))
		return nil, nil
	}

	// Cast the element's value to an Order
	order, ok := element.Value.(models.Order)
	if !ok {
		return nil, fmt.Errorf("failed to cast order for userId %q, symbol %q, and price %q", input.UserId, input.Symbol, input.Price)
	}

	return &order, nil
}
