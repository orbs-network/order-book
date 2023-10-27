package utils

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

type key int

const (
	userKey       key = iota
	paginationKey key = iota
)

// WithUserCtx returns a new context with the provided user value
func WithUserCtx(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUserCtx retrieves the user value from the context
func GetUserCtx(ctx context.Context) *models.User {
	if user, ok := ctx.Value(userKey).(*models.User); ok {
		return user
	}
	return nil
}

// WithPaginationCtx returns a new context with the provided pagination value
func WithPaginationCtx(ctx context.Context, pagination *Paginator) context.Context {
	return context.WithValue(ctx, paginationKey, pagination)
}

// GetPaginationCtx retrieves the pagination value from the context
func GetPaginationCtx(ctx context.Context) *Paginator {
	if pagination, ok := ctx.Value(paginationKey).(*Paginator); ok {
		return pagination
	}
	return nil
}
