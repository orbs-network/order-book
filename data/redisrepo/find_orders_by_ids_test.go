package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_FindOrdersByIds(t *testing.T) {

	orderId := uuid.MustParse("f6b9e7a0-9b9b-4e9f-8e1e-9e1b9d1c1b1b")

	t.Run("should return error if more than `MaxOrderIds` IDs are provided", func(t *testing.T) {
		ctx := context.Background()
		repo, _ := NewRedisRepository(nil)
		orders, err := repo.FindOrdersByIds(ctx, make([]uuid.UUID, MAX_ORDER_IDS+1), false)

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

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId}, false)

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

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId}, false)

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

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId}, false)

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

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{orderId}, false)

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

		orders, err := repo.FindOrdersByIds(ctx, []uuid.UUID{mocks.Order.Id}, false)

		assert.NoError(t, err)
		assert.Len(t, orders, 1)
		assert.Equal(t, mocks.Order.Id, orders[0].Id)
	})
}

func TestFindOrdersByIds_OnlyOpenFlag(t *testing.T) {
	type test struct {
		name              string
		orders            []models.Order
		onlyOpen          bool
		expectedID        uuid.UUID
		expectedNumOrders int
	}

	tests := []test{
		{
			name: "open order and pending order",
			orders: []models.Order{mocks.Order, {
				Id:          mocks.Order.Id,
				Size:        decimal.NewFromInt(123456),
				SizePending: decimal.NewFromInt(23456),
				Symbol:      mocks.Order.Symbol,
				Price:       mocks.Order.Price,
				Side:        mocks.Order.Side,
				Timestamp:   mocks.Order.Timestamp,
			}},
			onlyOpen:          true,
			expectedID:        mocks.Order.Id,
			expectedNumOrders: 2,
		},
		{
			name: "open order and filled order",
			orders: []models.Order{mocks.Order, {
				Id:         mocks.Order.Id,
				Size:       decimal.NewFromInt(1),
				SizeFilled: decimal.NewFromInt(1),
				Symbol:     mocks.Order.Symbol,
				Price:      mocks.Order.Price,
				Side:       mocks.Order.Side,
				Timestamp:  mocks.Order.Timestamp,
			}},
			onlyOpen:          true,
			expectedID:        mocks.Order.Id,
			expectedNumOrders: 1,
		},
		{
			name: "open order, pending order, filled order",
			orders: []models.Order{mocks.Order, {
				Id:          mocks.Order.Id,
				Size:        decimal.NewFromInt(123456),
				SizePending: decimal.NewFromInt(23456),
				Symbol:      mocks.Order.Symbol,
				Price:       mocks.Order.Price,
				Side:        mocks.Order.Side,
				Timestamp:   mocks.Order.Timestamp,
			}, {
				Id:         mocks.Order.Id,
				Size:       decimal.NewFromInt(1),
				SizeFilled: decimal.NewFromInt(1),
				Symbol:     mocks.Order.Symbol,
				Price:      mocks.Order.Price,
				Side:       mocks.Order.Side,
				Timestamp:  mocks.Order.Timestamp,
			}},
			onlyOpen:          true,
			expectedID:        mocks.Order.Id,
			expectedNumOrders: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			db, mock := redismock.NewClientMock()
			repo := &redisRepository{client: db}

			ids := make([]uuid.UUID, len(tc.orders))
			for i, order := range tc.orders {
				mock.ExpectHGetAll(CreateOrderIDKey(order.Id)).SetVal(order.OrderToMap())
				ids[i] = order.Id
			}

			orders, err := repo.FindOrdersByIds(ctx, ids, tc.onlyOpen)

			assert.NoError(t, err)
			assert.Len(t, orders, tc.expectedNumOrders)
			assert.Equal(t, tc.expectedID, orders[0].Id)
		})
	}

	t.Run("should return all orders if onlyOpen is false", func(t *testing.T) {
		ctx := context.Background()
		db, mock := redismock.NewClientMock()
		repo := &redisRepository{client: db}

		orders := []models.Order{{
			Id:         uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Size:       decimal.NewFromInt(123456),
			SizeFilled: decimal.NewFromInt(23456),
			Symbol:     mocks.Order.Symbol,
			Price:      mocks.Order.Price,
			Side:       mocks.Order.Side,
			Timestamp:  mocks.Order.Timestamp,
		},
			{
				Id:          uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				Size:        decimal.NewFromInt(123456),
				SizeFilled:  decimal.NewFromInt(23456),
				SizePending: decimal.NewFromInt(5555),
				Symbol:      mocks.Order.Symbol,
				Price:       mocks.Order.Price,
				Side:        mocks.Order.Side,
				Timestamp:   mocks.Order.Timestamp,
			},
			mocks.Order,
		}

		ids := make([]uuid.UUID, len(orders))
		for i, order := range orders {
			mock.ExpectHGetAll(CreateOrderIDKey(order.Id)).SetVal(order.OrderToMap())
			ids[i] = order.Id
		}

		orders, err := repo.FindOrdersByIds(ctx, ids, false)

		assert.NoError(t, err)
		assert.Len(t, orders, 3)
	})
}
