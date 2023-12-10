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
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var userId = uuid.MustParse("00000000-0000-0000-0000-000000000003")
var size, _ = decimal.NewFromString("10000324.123456789")

var order = models.Order{
	Id:     orderId,
	Price:  price,
	Size:   size,
	Symbol: symbol,
	Side:   models.BUY,
}

func TestRedisRepository_RemoveOrder(t *testing.T) {

	t.Run("fails to remove order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetErr(assert.AnError)

		err := repo.RemoveOrder(ctx, order)

		assert.ErrorContains(t, err, "failed to remove order")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// t.Run("succesfully removes order", func(t *testing.T) {
	// 	db, mock := redismock.NewClientMock()

	// 	repo := &redisRepository{
	// 		client: db,
	// 	}

	// 	buyPricesKey := CreateBuySidePricesKey(order.Symbol)
	// 	userOrdersKey := CreateUserOrdersKey(order.UserId)
	// 	clientOIdKey := CreateClientOIDKey(order.ClientOId)

	// 	mock.ExpectTxPipeline()
	// 	mock.ExpectZRem(buyPricesKey, order.Id.String()).SetVal(1)
	// 	mock.ExpectZRem(userOrdersKey, order.Id.String()).SetVal(1)
	// 	mock.ExpectDel(clientOIdKey, order.ClientOId.String()).SetVal(1)
	// 	mock.ExpectTxPipelineExec()

	// 	err := repo.RemoveOrder(ctx, order)

	// 	assert.NoError(t, err)
	// 	assert.NoError(t, mock.ExpectationsWereMet())
	// })
}
