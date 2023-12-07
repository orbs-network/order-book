package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/models"
)

type MockUserStore struct {
	User  *models.User
	Error error
}

func (m *MockUserStore) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	return *m.User, m.Error
}

func (m *MockUserStore) GetUserByApiKey(ctx context.Context, apiKey string) (*models.User, error) {
	return m.User, m.Error
}

func (m *MockUserStore) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	return m.User, m.Error
}

func (m *MockUserStore) UpdateUser(ctx context.Context, input storeuser.UpdateUserInput) error {
	return m.Error
}
