package mocks

import (
	"context"

	"github.com/orbs-network/order-book/service"
)

type MockBcClient struct {
	IsVerified bool
	Error      error
}

func (m *MockBcClient) VerifySignature(ctx context.Context, input service.VerifySignatureInput) (bool, error) {
	return m.IsVerified, m.Error
}
