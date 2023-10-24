package redisrepo

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreOrder(t *testing.T) {

	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

	t.Run("store order success (buy side) - should set user orders set, order ID hash, buy prices sorted set", func(t *testing.T) {
		var buyOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    symbol,
			Side:      models.BUY,
			Status:    models.STATUS_OPEN,
			Timestamp: timestamp,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectSAdd(CreateUserOrdersKey(buyOrder.UserId), buyOrder.Id.String()).SetVal(1)
		mock.ExpectHSet(CreateOrderIDKey(buyOrder.Id), buyOrder.OrderToMap()).SetVal(1)
		mock.ExpectSet(CreateClientOIDKey(buyOrder.ClientOId), buyOrder.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateBuySidePricesKey(buyOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: buyOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.StoreOrder(ctx, buyOrder)

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store order success (sell side) - should set user orders set, order ID hash, sell prices sorted set", func(t *testing.T) {
		var buyOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    symbol,
			Side:      models.SELL,
			Status:    models.STATUS_OPEN,
			Timestamp: timestamp,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectSAdd(CreateUserOrdersKey(buyOrder.UserId), buyOrder.Id.String()).SetVal(1)
		mock.ExpectHSet(CreateOrderIDKey(buyOrder.Id), buyOrder.OrderToMap()).SetVal(1)
		mock.ExpectSet(CreateClientOIDKey(buyOrder.ClientOId), buyOrder.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateSellSidePricesKey(buyOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: buyOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.StoreOrder(ctx, buyOrder)

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store order fail - should return error when transaction fails", func(t *testing.T) {

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectSAdd(CreateUserOrdersKey(order.UserId), order.Id.String()).SetErr(assert.AnError)
		mock.ExpectHSet(CreateOrderIDKey(buyOrder.Id), order.OrderToMap()).SetErr(assert.AnError)
		mock.ExpectSet(CreateClientOIDKey(buyOrder.ClientOId), buyOrder.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateSellSidePricesKey(order.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: order.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec().SetErr(fmt.Errorf("transaction failed"))

		err := repo.StoreOrder(ctx, order)

		assert.ErrorContains(t, err, fmt.Sprintf("transaction failed. Reason: %v", assert.AnError), "should return error")
	})

}
