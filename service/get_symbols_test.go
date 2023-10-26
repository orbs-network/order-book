package service_test

import (
	"context"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

func TestService_GetSymbols(t *testing.T) {
	ctx := context.Background()
	svc, _ := service.New(&mocks.MockOrderBookStore{})

	symbols, err := svc.GetSymbols(ctx)

	assert.Greater(t, len(symbols), 20)
	assert.NoError(t, err)
}
