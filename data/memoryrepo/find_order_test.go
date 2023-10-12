package memoryrepo

import (
	"container/list"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_FindOrderByUserAndPrice(t *testing.T) {
	repo, _ := NewMemoryRepository()
	ctx := context.Background()

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

	err := repo.StoreOrder(ctx, order)
	assert.NoError(t, err)

	t.Run("order found", func(t *testing.T) {
		foundOrder, err := repo.FindOrder(ctx, models.FindOrderInput{
			UserId: userId,
			Price:  price,
			Symbol: symbol,
		})
		assert.NoError(t, err)
		assert.Equal(t, order, *foundOrder)
	})

	t.Run("when user does not exist", func(t *testing.T) {
		nonExistentUserId := uuid.New()
		foundOrder, err := repo.FindOrder(ctx, models.FindOrderInput{
			UserId: nonExistentUserId,
			Price:  price,
			Symbol: symbol,
		})
		assert.NoError(t, err, "user does not exist error should be nil")
		assert.Nil(t, foundOrder, "order should not be found")
	})

	t.Run("when symbol does not exist", func(t *testing.T) {
		nonExistentSymbol, _ := models.StrToSymbol("BTC-ETH")
		foundOrder, err := repo.FindOrder(ctx, models.FindOrderInput{
			UserId: userId,
			Price:  price,
			Symbol: nonExistentSymbol,
		})
		assert.NoError(t, err, "symbol does not exist error should be nil")
		assert.Nil(t, foundOrder, "order should not be found")
	})

	t.Run("when price does not exist", func(t *testing.T) {
		nonExistentPrice := decimal.NewFromFloat(20.0)
		foundOrder, err := repo.FindOrder(ctx, models.FindOrderInput{
			UserId: userId,
			Price:  nonExistentPrice,
			Symbol: symbol,
		})
		assert.NoError(t, err, "price does not exist error should be nil")
		assert.Nil(t, foundOrder)
	})

	t.Run("cast error", func(t *testing.T) {
		invalidOrder := "invalid order"
		repo.userOrders[userId.String()] = map[models.Symbol]map[string]*list.Element{
			symbol: {
				price.StringFixed(models.STR_PRECISION): list.New().PushBack(invalidOrder),
			},
		}
		foundOrder, err := repo.FindOrder(ctx, models.FindOrderInput{
			UserId: userId,
			Price:  price,
			Symbol: symbol,
		})
		assert.EqualError(t, err, fmt.Sprintf("failed to cast order for userId %q, symbol %q, and price %q", userId, symbol, price))
		assert.Nil(t, foundOrder)
	})
}
