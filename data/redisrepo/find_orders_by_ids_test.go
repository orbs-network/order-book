package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_FindOrdersByIds(t *testing.T) {

	orderId := uuid.MustParse("f6b9e7a0-9b9b-4e9f-8e1e-9e1b9d1c1b1b")

	t.Run("should return error if more than `MaxOrderIds` IDs are provided", func(t *testing.T) {
		ctx := context.Background()
		repo, _ := NewRedisRepository(nil)
		orders, err := repo.FindOrdersByIds(ctx, make([]uuid.UUID, MAX_ORDER_IDS+1))

		assert.ErrorContains(t, err, "exceeded maximum number of IDs")
		assert.Nil(t, orders)
	})

	t.Run("should return error if order is not found", func(t *testing.T) {

		ctx := context.Background()

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectHGetAll(CreateOrderIDKey(orderId)).SetVal(map[string]string{})

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId})

		assert.ErrorContains(t, err, "order not found but was expected to exist")
		assert.Nil(t, orders)
	})

	t.Run("should return error if order was expected to exist but was not found", func(t *testing.T) {

		ctx := context.Background()

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectHGetAll(CreateOrderIDKey(orderId)).SetVal(map[string]string{})

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId})

		assert.ErrorContains(t, err, "order not found but was expected to exist")
		assert.Nil(t, orders)
	})

	t.Run("should return error if order could not be mapped", func(t *testing.T) {

		ctx := context.Background()

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectHGetAll(CreateOrderIDKey(orderId)).SetVal(map[string]string{"id": "invalid"})

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId})

		assert.ErrorContains(t, err, "could not map order")
		assert.Nil(t, orders)
	})

	t.Run("should return error if pipeline could not be executed", func(t *testing.T) {

		ctx := context.Background()

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectHGetAll(CreateOrderIDKey(orderId)).SetVal(map[string]string{"id": "invalid"})

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId})

		assert.ErrorContains(t, err, "could not map order")
		assert.Nil(t, orders)
	})

	t.Run("should return orders if found", func(t *testing.T) {

		ctx := context.Background()

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectHGetAll(CreateOrderIDKey(mocks.Order.Id)).SetVal(mocks.Order.OrderToMap())

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{mocks.Order.Id})

		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, mocks.Order.Id, orders[0].Id)
	})
}
