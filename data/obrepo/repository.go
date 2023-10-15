package obrepo

import (
	"container/list"
	"sync"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// Represents a collection of orders at a given price for a given symbol (e.g. USDC-ETH).
type ordersAtPrice struct {
	Price  decimal.Decimal
	Orders *list.List
	Index  int // position in heap
}

// Represents a single order book for a given symbol (e.g. USDC-ETH).
type orderBook struct {
	sellOrders    sellHeap
	buyOrders     buyHeap
	orderIDMap    map[uuid.UUID]models.Order
	userOrdersMap map[uuid.UUID][]models.Order
	priceCache    map[decimal.Decimal][]models.Order
	mu            sync.RWMutex
}

// Represents a collection of order books.
type orderBookManager struct {
	orderBooks map[models.Symbol]*orderBook
	mu         sync.RWMutex
}

type OrderBookManager interface {
	GetOrderBook(symbol models.Symbol) (*orderBook, error)
	CreateOrderBook(symbol models.Symbol) (*orderBook, error)
}

// Creates a new order book for a given symbol (e.g. USDC-ETH).
func NewOrderBook() *orderBook {
	return &orderBook{
		orderIDMap:    make(map[uuid.UUID]models.Order),
		userOrdersMap: make(map[uuid.UUID][]models.Order),
		priceCache:    make(map[decimal.Decimal][]models.Order),
	}
}

// Creates a manager for multiple order book instances.
func NewOrderBookManager() *orderBookManager {
	return &orderBookManager{
		orderBooks: make(map[models.Symbol]*orderBook),
	}
}
