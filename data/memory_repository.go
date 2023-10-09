// Package data responsible for handling requests to the underlying repository.
//
// data returned from this layer should be mapped into a known type within the service
// to avoid leaking internal implementation types.

package data

import "github.com/orbs-network/order-book/models"

// MemoryRespository contains methods for interacting with an in memory repository.
type memoryRespository struct {
	Orders map[string]models.Order
}

// NewMemoryRespository creates a new memory repository.
func NewMemoryRespository() (*memoryRespository, error) {
	return &memoryRespository{
		Orders: make(map[string]models.Order),
	}, nil
}
