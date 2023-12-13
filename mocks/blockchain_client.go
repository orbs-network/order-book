package mocks

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
)

type MockBcClient struct {
	IsVerified bool
	Error      error
	Tx         models.Tx
}

func (m *MockBcClient) VerifySignature(ctx context.Context, input service.VerifySignatureInput) (bool, error) {
	return m.IsVerified, m.Error
}

func (m *MockBcClient) CheckPendingTxs(ctx context.Context) error {
	return m.Error
}
