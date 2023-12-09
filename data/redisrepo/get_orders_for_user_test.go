package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_GetOrdersForUser(t *testing.T) {

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	var order = models.Order{
		Id:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		UserId: userId,
		Symbol: "BTC-ETH",
		Side:   models.BUY,
	}

	t.Run("should get orders for user", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		ctx := mocks.AddPaginationToCtx(1, 10)

		key := CreateUserOpenOrdersKey(userId)

		mock.ExpectZCard(key).SetVal(1)
		mock.ExpectZRange(key, int64(0), int64(10)).SetVal([]string{"00000000-0000-0000-0000-000000000001"})
		mock.ExpectHGetAll(CreateOrderIDKey(order.Id)).SetVal(order.OrderToMap())

		orders, totalOrders, err := repo.GetOrdersForUser(ctx, userId, false)

		assert.Equal(t, orders[0].Id, order.Id)
		assert.Equal(t, orders[0].UserId, order.UserId)
		assert.Len(t, orders, 1)
		assert.Equal(t, 1, totalOrders)
		assert.NoError(t, err)
	})

	t.Run("should return error if failed to get total count of orders for user", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		ctx := mocks.AddPaginationToCtx(1, 10)

		key := CreateUserOpenOrdersKey(userId)

		mock.ExpectZCard(key).SetErr(assert.AnError)

		orders, totalOrders, err := repo.GetOrdersForUser(ctx, userId, false)

		assert.Equal(t, orders, []models.Order{})
		assert.Equal(t, totalOrders, 0)
		assert.ErrorContains(t, err, "failed to get total count of orders for user")
	})

	t.Run("should return error if failed to get orders for user", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		ctx := mocks.AddPaginationToCtx(1, 10)

		key := CreateUserOpenOrdersKey(userId)

		mock.ExpectZCard(key).SetVal(1)
		mock.ExpectZRange(key, int64(0), int64(10)).SetErr(assert.AnError)

		orders, totalOrders, err := repo.GetOrdersForUser(ctx, userId, false)

		assert.Equal(t, orders, []models.Order{})
		assert.Equal(t, totalOrders, 0)
		assert.ErrorContains(t, err, "failed to get orders for user")
	})

	t.Run("should return error if failed to parse order id", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		ctx := mocks.AddPaginationToCtx(1, 10)

		key := CreateUserOpenOrdersKey(userId)

		mock.ExpectZCard(key).SetVal(1)
		mock.ExpectZRange(key, int64(0), int64(10)).SetVal([]string{"bad-uuid"})

		orders, totalOrders, err := repo.GetOrdersForUser(ctx, userId, false)

		assert.Equal(t, orders, []models.Order{})
		assert.Equal(t, totalOrders, 0)
		assert.ErrorContains(t, err, "failed to parse order id")
	})

	t.Run("should return error if failed to get order", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		ctx := mocks.AddPaginationToCtx(1, 10)

		key := CreateUserOpenOrdersKey(userId)

		mock.ExpectZCard(key).SetVal(1)
		mock.ExpectZRange(key, int64(0), int64(10)).SetVal([]string{"00000000-0000-0000-0000-000000000001"})
		mock.ExpectHGetAll(CreateOrderIDKey(order.Id)).SetErr(assert.AnError)

		orders, totalOrders, err := repo.GetOrdersForUser(ctx, userId, false)

		assert.Equal(t, orders, []models.Order{})
		assert.Equal(t, totalOrders, 0)
		assert.ErrorContains(t, err, "failed to get order")
	})

}
