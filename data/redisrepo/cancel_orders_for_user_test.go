package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_CancelAllOrdersForUser(t *testing.T) {
	ctx := context.Background()

	orderTwo := mocks.Order
	orderTwo.Side = models.SELL

	t.Run("should remove single order for user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(mocks.OrderId)).SetVal(mocks.Order.OrderToMap())
		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateClientOIDKey(mocks.ClientOId)).SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(mocks.Symbol), mocks.OrderId.String()).SetVal(1)
		mock.ExpectDel(CreateOrderIDKey(mocks.OrderId)).SetVal(1)
		mock.ExpectDel(CreateUserOrdersKey(mocks.UserId)).SetVal(1)
		mock.ExpectTxPipelineExec()

		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Equal(t, mocks.OrderId, orderIds[0])
		assert.Len(t, orderIds, 1)
		assert.NoError(t, err)
	})

	t.Run("should remove multiple orders for user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String(), orderTwo.Id.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(mocks.OrderId)).SetVal(mocks.Order.OrderToMap())
		mock.ExpectHGetAll(CreateOrderIDKey(orderTwo.Id)).SetVal(orderTwo.OrderToMap())
		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateClientOIDKey(mocks.ClientOId)).SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(mocks.Symbol), mocks.OrderId.String()).SetVal(1)
		mock.ExpectDel(CreateOrderIDKey(mocks.OrderId)).SetVal(1)
		mock.ExpectDel(CreateClientOIDKey(orderTwo.ClientOId)).SetVal(1)
		mock.ExpectZRem(CreateSellSidePricesKey(orderTwo.Symbol), orderTwo.Id.String()).SetVal(1)
		mock.ExpectDel(CreateOrderIDKey(orderTwo.Id)).SetVal(1)
		mock.ExpectDel(CreateUserOrdersKey(mocks.UserId)).SetVal(1)
		mock.ExpectTxPipelineExec()

		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Equal(t, mocks.OrderId, orderIds[0])
		assert.Equal(t, orderTwo.Id, orderIds[1])
		assert.Len(t, orderIds, 2)
		assert.NoError(t, err)
	})

	t.Run("should exit with error if failed to get order IDs for user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetErr(assert.AnError)
		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)

		assert.Empty(t, orderIds)
		assert.ErrorContains(t, err, "failed to get order IDs for user")
	})

	t.Run("should return `ErrNoOrdersFound` error if no orders are found for the user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{})

		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Empty(t, orderIds)
		assert.ErrorIs(t, err, models.ErrNoOrdersFound)
	})

	t.Run("should exit with error if failed to parse order ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{"invalid"})
		_, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "failed to parse order ID: invalid UUID length: 7")
	})

	t.Run("should immediately return error if one order is not found by ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String(), orderTwo.Id.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(orderTwo.Id)).SetErr(models.ErrOrderNotFound) // order not found - break out of loop iteration

		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Empty(t, orderIds)
		assert.ErrorContains(t, err, "failed to find orders by IDs")
	})

	t.Run("an error should be returned if transaction failed", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(mocks.OrderId)).SetVal(mocks.Order.OrderToMap())
		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateClientOIDKey(mocks.ClientOId)).SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(mocks.Symbol), mocks.OrderId.String()).SetVal(1)
		mock.ExpectDel(CreateOrderIDKey(mocks.OrderId)).SetVal(1)
		mock.ExpectDel(CreateUserOrdersKey(mocks.UserId)).SetErr(assert.AnError)
		mock.ExpectTxPipelineExec().SetErr(assert.AnError)

		orderIds, err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.Empty(t, orderIds)
		assert.ErrorContains(t, err, "failed to remove user's orders in Redis. Reason: assert.AnError")
	})

}
