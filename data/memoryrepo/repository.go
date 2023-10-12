package memoryrepo

import (
	"container/list"
	"sync"

	"github.com/shopspring/decimal"
)

type inMemoryRepository struct {
	// Map of sell orders, keyed by price
	sellOrders map[decimal.Decimal]*ordersAtPrice
	// Map of order IDs to their location in the sellOrders map
	orderLocations map[string]*list.Element
	// Map of orders, keyed by user id and price. Prevents duplicate orders
	userOrders map[string]map[decimal.Decimal]*list.Element
	// Mutex to protect concurrent access to the above maps
	mu sync.RWMutex
}

type ordersAtPrice struct {
	// Double linked list of orders at this price
	List *list.List
	// Sum of all orders at this price
	Sum decimal.Decimal
}

func NewMemoryRepository() (*inMemoryRepository, error) {
	return &inMemoryRepository{
		sellOrders:     make(map[decimal.Decimal]*ordersAtPrice),
		orderLocations: make(map[string]*list.Element),
		userOrders:     make(map[string]map[decimal.Decimal]*list.Element),
	}, nil
}
