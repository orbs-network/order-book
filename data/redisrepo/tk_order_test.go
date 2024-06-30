package redisrepo

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Init Redis repository with a mock client and a mock tx
func setupTest() (context.Context, redismock.ClientMock, *redisRepository) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	repo := &redisRepository{
		client: db,
		txMap:  make(map[uint]redis.Pipeliner),
	}

	tx := db.TxPipeline()
	txid := uint(1)

	repo.txMap[txid] = tx

	return ctx, mock, repo
}

func TestRedisRepository_TxStartEndPerform(t *testing.T) {

	t.Run("txStart initializes a transaction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client:  db,
			txMap:   make(map[uint]redis.Pipeliner),
			ixIndex: 0,
		}

		mock.ExpectTxPipeline()

		txid := repo.txStart(context.Background())

		assert.Equal(t, uint(1), txid)
		assert.Contains(t, repo.txMap, txid)
	})

	t.Run("txEnd executes and clears transaction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		pipeline := db.TxPipeline()

		repo := &redisRepository{
			client: db,
			txMap:  map[uint]redis.Pipeliner{1: pipeline},
		}

		err := repo.txEnd(context.Background(), 1)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.NotContains(t, repo.txMap, 1)
	})

	t.Run("PerformTx executes action within a transaction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client:  db,
			txMap:   make(map[uint]redis.Pipeliner),
			ixIndex: 0,
		}

		actionCalled := false
		testAction := func(txid uint) error {
			actionCalled = true
			return nil
		}

		err := repo.PerformTx(context.Background(), testAction)

		assert.NoError(t, err)
		assert.True(t, actionCalled)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRedisRepository_TxModifyOrder(t *testing.T) {

	ctx, mock, repo := setupTest()

	t.Run("successfully adds order", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(mocks.Order.Id), mocks.Order.OrderToMap()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyOrder(ctx, txid, models.Add, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully updates order", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(mocks.Order.Id), mocks.Order.OrderToMap()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyOrder(ctx, txid, models.Update, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully deletes order", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateOrderIDKey(mocks.Order.Id)).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyOrder(ctx, txid, models.Remove, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unsupported operation", func(t *testing.T) {

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyOrder(ctx, txid, 99, mocks.Order)
		})

		assert.ErrorIs(t, err, models.ErrUnsupportedOperation)
	})
}

func TestRedisRepository_TxModifyPrices(t *testing.T) {

	ctx, mock, repo := setupTest()

	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

	t.Run("successfully adds to buy side prices key", func(t *testing.T) {

		var buyOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    test_symbol,
			Side:      models.BUY,
			Timestamp: timestamp,
		}

		mock.ExpectTxPipeline()
		mock.ExpectZAdd(CreateBuySidePricesKey(buyOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: buyOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyPrices(ctx, txid, models.Add, buyOrder)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully adds to sell side prices key", func(t *testing.T) {

		var sellOrder = models.Order{
			Id:        orderId,
			ClientOId: clientOId,
			Price:     price,
			Size:      size,
			Symbol:    test_symbol,
			Side:      models.SELL,
			Timestamp: timestamp,
		}

		mock.ExpectTxPipeline()
		mock.ExpectZAdd(CreateSellSidePricesKey(sellOrder.Symbol), redis.Z{
			Score:  10.0016969392,
			Member: sellOrder.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyPrices(ctx, txid, models.Add, sellOrder)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully deletes price", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectZRem(CreateBuySidePricesKey(mocks.Order.Symbol), mocks.Order.Id.String()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyPrices(ctx, txid, models.Remove, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unsupported operation", func(t *testing.T) {

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyPrices(ctx, txid, 99, mocks.Order)
		})

		assert.ErrorIs(t, err, models.ErrUnsupportedOperation)
	})
}

func TestRedisRepository_TxModifyClientOId(t *testing.T) {

	ctx, mock, repo := setupTest()

	t.Run("successfully adds clientOID", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectSet(CreateClientOIDKey(mocks.Order.ClientOId), mocks.Order.Id.String(), 0).SetVal("OK")
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyClientOId(ctx, txid, models.Add, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully deletes clientOID", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectDel(CreateClientOIDKey(mocks.Order.ClientOId)).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyClientOId(ctx, txid, models.Remove, mocks.Order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unsupported operation", func(t *testing.T) {

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyClientOId(ctx, txid, 99, mocks.Order)
		})

		assert.ErrorIs(t, err, models.ErrUnsupportedOperation)
	})
}

func TestRedisRepository_TxModifyUserOpenOrders(t *testing.T) {

	ctx, mock, repo := setupTest()
	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	order := models.Order{
		Timestamp: timestamp,
		Symbol:    test_symbol,
	}

	t.Run("successfully adds user open order", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectZAdd(CreateUserOpenOrdersKey(order.UserId, order.Symbol), redis.Z{
			Score:  float64(timestamp.UnixNano()),
			Member: order.Id.String(),
		}).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyUserOpenOrders(ctx, txid, models.Add, order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("successfully removes user open order", func(t *testing.T) {

		mock.ExpectTxPipeline()
		mock.ExpectZRem(CreateUserOpenOrdersKey(test_order.UserId, test_symbol), test_order.Id.String()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyUserOpenOrders(ctx, txid, models.Remove, test_order)
		})

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("unsupported operation", func(t *testing.T) {

		err := repo.PerformTx(ctx, func(txid uint) error {
			return repo.TxModifyUserOpenOrders(ctx, txid, 99, mocks.Order)
		})

		assert.ErrorIs(t, err, models.ErrUnsupportedOperation)
	})
}
