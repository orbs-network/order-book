package redisrepo

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreFilledOrders(t *testing.T) {
	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

	t.Run("store filled order success (buy side) - should remove order from user open orders set, add order to user filled orders set, remove order from buy prices sorted set", func(t *testing.T) {
		var buyOrder = models.Order{
			Id:          orderId,
			ClientOId:   clientOId,
			Price:       price,
			Size:        size,
			Symbol:      symbol,
			Side:        models.BUY,
			Timestamp:   timestamp,
			SizePending: decimal.Zero,
			SizeFilled:  size,
			UserId:      userId,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectZRem(CreateUserOpenOrdersKey(buyOrder.UserId, symbol), buyOrder.Id.String()).SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(buyOrder.Symbol), buyOrder.Id.String()).SetVal(1)
		mock.ExpectHSet(CreateOrderIDKey(buyOrder.Id), buyOrder.OrderToMap()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.StoreFilledOrders(ctx, []models.Order{buyOrder})

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store filled order success (sell side) - should remove order from user open orders set, add order to user filled orders set, remove order from sell prices sorted set", func(t *testing.T) {
		var sellOrder = models.Order{
			Id:          orderId,
			ClientOId:   clientOId,
			Price:       price,
			Size:        size,
			Symbol:      symbol,
			Side:        models.SELL,
			Timestamp:   timestamp,
			SizePending: decimal.Zero,
			SizeFilled:  size,
			UserId:      userId,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectZRem(CreateUserOpenOrdersKey(sellOrder.UserId, symbol), sellOrder.Id.String()).SetVal(1)
		mock.ExpectZRem(CreateSellSidePricesKey(sellOrder.Symbol), sellOrder.Id.String()).SetVal(1)
		mock.ExpectHSet(CreateOrderIDKey(sellOrder.Id), sellOrder.OrderToMap()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.StoreFilledOrders(ctx, []models.Order{sellOrder})

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store filled order unexpected error - should return error", func(t *testing.T) {
		var sellOrder = models.Order{
			UserId: userId,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectZRem(CreateUserOpenOrdersKey(sellOrder.UserId, symbol), sellOrder.Id.String()).SetErr(assert.AnError)

		err := repo.StoreFilledOrders(ctx, []models.Order{sellOrder})

		assert.ErrorContains(t, err, "failed to store filled orders in Redis")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
