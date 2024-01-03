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
	svc, _ := service.New(&mocks.MockOrderBookStore{}, &mocks.MockBcClient{})

	symbols, err := svc.GetSymbols(ctx)

	assert.GreaterOrEqual(t, len(symbols), 1)
	assert.NoError(t, err)
}
