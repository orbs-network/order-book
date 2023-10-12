// Package service contains the business logic for the application.

package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type Store interface {
	StoreOrder(order models.Order) error
	RemoveOrder(orderId uuid.UUID) error
	GetOrdersAtPrice(price decimal.Decimal) []models.Order
	GetAllPrices() []decimal.Decimal
}

// Service contains methods that implement the business logic for the application.
type Service struct {
	store Store
}

// New creates a new Service with injected dependencies.
func New(store Store) (*Service, error) {
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	return &Service{store: store}, nil
}
