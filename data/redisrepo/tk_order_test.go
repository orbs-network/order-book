package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_TxModifyOrder(t *testing.T) {

	t.Run("txStart initializes a transaction", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client:  db,
			txMap:   make(map[uint]redis.Pipeliner),
			ixIndex: 0,
		}

		// Set expectations
		mock.ExpectTxPipeline()

		txid, err := repo.txStart(context.Background())

		assert.NoError(t, err)
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

		repo.txEnd(context.Background(), 1)

		// Check that the transaction was attempted to be executed
		assert.NoError(t, mock.ExpectationsWereMet())

		// Ensure the txMap no longer contains the transaction
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

func TestRedisRepository_TxDeleteClientOID(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()

	repo := &redisRepository{
		client: db,
		txMap:  make(map[uint]redis.Pipeliner),
	}

	tx := db.TxPipeline()
	txid := uint(1)

	repo.txMap[txid] = tx

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
}
