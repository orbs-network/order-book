package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

type MockBcClient struct {
	IsVerified bool
	Error      error
	Tx         models.Tx
}

func (m *MockBcClient) CheckPendingTxs(ctx context.Context) error {
	return m.Error
}
