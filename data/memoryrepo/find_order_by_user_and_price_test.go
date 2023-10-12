package memoryrepo

import (
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_FindOrderByUserAndPrice(t *testing.T) {
	repo, _ := NewMemoryRepository()

	symbol, _ := models.StrToSymbol("USDC-ETH")

	userId := uuid.New()
	price := decimal.NewFromFloat(10.0)
	order := models.Order{
		Id:     uuid.New(),
		UserId: userId,
		Price:  price,
		Size:   decimal.NewFromFloat(1.0),
		Symbol: symbol,
	}

	err := repo.StoreOrder(order)
	assert.NoError(t, err)

	t.Run("when user exists and price exists", func(t *testing.T) {
		foundOrder, err := repo.FindOrderByUserAndPrice(userId, price)
		assert.NoError(t, err)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		nonExistentUserId := uuid.New()
		foundOrder, err := repo.FindOrderByUserAndPrice(nonExistentUserId, price)
		assert.NoError(t, err)
		assert.Nil(t, foundOrder)
	})

	t.Run("when price does not exist", func(t *testing.T) {
		nonExistentPrice := decimal.NewFromFloat(20.0)
		foundOrder, err := repo.FindOrderByUserAndPrice(userId, nonExistentPrice)
		assert.NoError(t, err)
		assert.Nil(t, foundOrder)
	})
}
