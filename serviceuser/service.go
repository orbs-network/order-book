package serviceuser

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/models"
)

type UserService interface {
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	GetUserByApiKey(ctx context.Context, apiKey string) (*models.User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) error
}

type Service struct {
	userStore storeuser.UserStore
}

func New(userStore storeuser.UserStore) (UserService, error) {
	if userStore == nil {
		return nil, errors.New("userStore is nil")
	}

	return &Service{userStore: userStore}, nil
}
