package memoryrepo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *inMemoryRepository) FindOrderByUserAndPrice(userId uuid.UUID, price decimal.Decimal) (*models.Order, error) {

	userOrders, userExists := r.userOrders[userId.String()]
	if !userExists {
		fmt.Println("No user found for that order")
		return nil, nil
	}

	element, priceExists := userOrders[price.StringFixed(models.STR_PRECISION)]

	if !priceExists {
		fmt.Println("No price found for that order")
		return nil, nil
	}

	order, ok := element.Value.(models.Order)
	if !ok {
		return nil, fmt.Errorf("failed to cast order for userId %q and price %q", userId, price.String())
	}

	return &order, nil
}
