// Package service contains the business logic for the application.

package service

import (
	"context"
	"errors"

	"github.com/orbs-network/order-book/models"
)

type Store interface {
	AddOrder(ctx context.Context, order models.Order) (models.Order, error)
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
