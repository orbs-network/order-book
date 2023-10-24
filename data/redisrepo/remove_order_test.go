package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var orderId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var size, _ = decimal.NewFromString("10000324.123456789")

var order = models.Order{
	Id:     orderId,
	Price:  price,
	Size:   size,
	Symbol: symbol,
	Side:   models.BUY,
	Status: models.STATUS_OPEN,
}

func TestRedisRepository_RemoveOrder(t *testing.T) {

	t.Run("only open orders can be removed", func(t *testing.T) {
		db, _ := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		err := repo.RemoveOrder(ctx, models.Order{Status: models.STATUS_PENDING})

		assert.ErrorIs(t, err, models.ErrOrderNotOpen)
	})

	t.Run("fails to remove order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetErr(assert.AnError)

		err := repo.RemoveOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrTransactionFailed)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("succesfully removes order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		userOrdersKey := CreateUserOrdersKey(order.UserId)
		orderIDKey := CreateOrderIDKey(order.Id)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetVal(1)
		mock.ExpectSRem(userOrdersKey, order.Id.String()).SetVal(1)
		mock.ExpectHSet(orderIDKey, "status", models.STATUS_CANCELLED.String()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.RemoveOrder(ctx, order)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}