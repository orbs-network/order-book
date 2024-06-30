package redisrepo

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreOpenOrder(t *testing.T) {

	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	txMap := make(map[uint]redis.Pipeliner)
	t.Run("store open order success (buy side) - should set user orders set, order ID hash, buy prices sorted set", func(t *testing.T) {
		var buyOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    test_symbol,
			Side:      models.BUY,
			Timestamp: timestamp,
		}

		db, mock := redismock.NewClientMock()
		repo := &redisRepository{
			client: db,
			txMap:  txMap,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(buyOrder.Id), buyOrder.OrderToMap()).SetVal(1)
		mock.ExpectSet(CreateClientOIDKey(buyOrder.ClientOId), buyOrder.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateBuySidePricesKey(buyOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: buyOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectZAdd(CreateUserOpenOrdersKey(buyOrder.UserId, buyOrder.Symbol), redis.Z{
			Score:  float64(timestamp.UnixNano()),
			Member: buyOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectSetNX(Order2MakerTokenTrackKey(buyOrder), -1, 0).SetVal(true)
		mock.ExpectTxPipelineExec()

		err := repo.StoreOpenOrder(ctx, buyOrder)

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store open order success (sell side) - should set user orders set, order ID hash, sell prices sorted set", func(t *testing.T) {
		var sellOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    test_symbol,
			Side:      models.SELL,
			Timestamp: timestamp,
		}

		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
			txMap:  txMap,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(sellOrder.Id), sellOrder.OrderToMap()).SetVal(1)
		mock.ExpectSet(CreateClientOIDKey(sellOrder.ClientOId), sellOrder.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateSellSidePricesKey(sellOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: sellOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectZAdd(CreateUserOpenOrdersKey(sellOrder.UserId, sellOrder.Symbol), redis.Z{
			Score:  float64(timestamp.UnixNano()),
			Member: sellOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectSetNX(Order2MakerTokenTrackKey(sellOrder), -1, 0).SetVal(true)
		mock.ExpectTxPipelineExec()

		err := repo.StoreOpenOrder(ctx, sellOrder)

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store open order fail - should return error when transaction fails", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
			txMap:  txMap,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(test_order.Id), test_order.OrderToMap()).SetErr(assert.AnError)
		mock.ExpectSet(CreateClientOIDKey(test_order.ClientOId), test_order.Id.String(), 0).SetVal("OK")
		mock.ExpectZAdd(CreateSellSidePricesKey(test_order.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: test_order.Id.String(),
		}).SetVal(1)
		mock.ExpectZAdd(CreateUserOpenOrdersKey(test_order.UserId, test_order.Symbol), redis.Z{
			Score:  float64(test_order.Timestamp.UnixNano()),
			Member: test_order.Id.String(),
		}).SetErr(assert.AnError)
		mock.ExpectExists(Order2MakerTokenTrackKey(test_order)).SetVal(0)
		mock.ExpectSetNX(Order2MakerTokenTrackKey(test_order), -1, 0).SetVal(true)

		err := repo.StoreOpenOrder(ctx, test_order)

		assert.ErrorContains(t, err, "PerformTx txEnd commit failed", "should return error")
	})

}
