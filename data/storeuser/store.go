package storeuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
)

type UpdateUserInput struct {
	UserId uuid.UUID
	PubKey string
	ApiKey string
}

type UserStore interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByApiKey(ctx context.Context, apiKey string) (*models.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) error
	// TODO: how should we handle removing users? What happens to their orders?
}
