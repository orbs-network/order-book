package utils

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

type key int

const (
	userKey key = iota
)

// WithUser returns a new context with the provided user value
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser retrieves the user value from the context
func GetUser(ctx context.Context) *models.User {
	if user, ok := ctx.Value(userKey).(*models.User); ok {
		return user
	}
	return nil
}
