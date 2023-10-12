package memoryrepo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_FindOrderById(t *testing.T) {
	repo, _ := NewMemoryRepository()
	ctx := context.Background()

	symbol, _ := models.StrToSymbol("USDC-ETH")
	orderId := uuid.New()

	userId := uuid.New()
	price := decimal.NewFromFloat(10.0)
	order := models.Order{
		Id:     orderId,
		UserId: userId,
		Price:  price,
		Size:   decimal.NewFromFloat(1.0),
		Symbol: symbol,
	}

	err := repo.StoreOrder(ctx, order)
	assert.NoError(t, err)

	t.Run("when order exists", func(t *testing.T) {
		foundOrder, err := repo.FindOrderById(ctx, orderId)
		assert.NoError(t, err)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("when order does not exist", func(t *testing.T) {
		nonExistentOrderId := uuid.New()
		foundOrder, err := repo.FindOrderById(ctx, nonExistentOrderId)
		assert.NoError(t, err, "order does not exist error should be nil")
		assert.Nil(t, foundOrder)
	})

}
