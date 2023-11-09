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

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
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

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.NoError(t, err)
	})

	t.Run("should exit with error if failed to get order IDs for user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetErr(assert.AnError)
		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "failed to fetch user order IDs. Reason: assert.AnError")
	})

	t.Run("should exit without error if no orders found for user", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{})

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.NoError(t, err)
	})

	t.Run("should exit with error if failed to parse order ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{"invalid"})
		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "failed to parse order ID: invalid UUID length: 7")
	})

	// TODO: confirm whether this is the behaviour we want
	t.Run("should continue removing other orders if one order is not found by ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String(), orderTwo.Id.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(orderTwo.Id)).SetErr(models.ErrOrderNotFound) // order not found - break out of loop iteration
		mock.ExpectHGetAll(CreateOrderIDKey(orderTwo.Id)).SetVal(mocks.Order.OrderToMap())
		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateClientOIDKey(orderTwo.ClientOId)).SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(orderTwo.Symbol), orderTwo.Id.String()).SetVal(1)
		mock.ExpectDel(CreateOrderIDKey(orderTwo.Id)).SetVal(1)
		mock.ExpectDel(CreateUserOrdersKey(orderTwo.UserId)).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.NoError(t, err)
	})

	t.Run("should exit with error if failed to get order by ID", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectZRange(CreateUserOrdersKey(mocks.UserId), 0, -1).SetVal([]string{mocks.OrderId.String()})
		mock.ExpectHGetAll(CreateOrderIDKey(mocks.OrderId)).SetErr(assert.AnError)

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "unexpected error finding order by ID: assert.AnError")
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

		err := repo.CancelOrdersForUser(ctx, mocks.UserId)
		assert.ErrorContains(t, err, "failed to remove user's orders in Redis. Reason: assert.AnError")
	})

}
