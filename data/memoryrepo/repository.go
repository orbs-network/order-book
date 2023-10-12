package memoryrepo

import (
	"container/list"
	"sync"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type inMemoryRepository struct {
	// symbol -> price -> (orders and sum)
	sellOrders map[models.Symbol]map[string]*ordersAtPrice
	// orderId -> (order)
	orderLocations map[string]*list.Element
	// userId -> symbol -> price -> (order) . Allows multiple orders per symbol per price
	userOrders map[string]map[models.Symbol]map[string]*list.Element
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
		sellOrders:     make(map[models.Symbol]map[string]*ordersAtPrice),
		orderLocations: make(map[string]*list.Element),
		userOrders:     make(map[string]map[models.Symbol]map[string]*list.Element),
	}, nil
}
