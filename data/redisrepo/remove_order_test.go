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

func TestRedisRepository_RemoveOrder(t *testing.T) {

	ctx := context.Background()
	orderId := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	size, _ := decimal.NewFromString("10000324.123456789")

	t.Run("only open orders can be removed", func(t *testing.T) {
		db, _ := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		err := repo.RemoveOrder(ctx, models.Order{Status: models.STATUS_PENDING})

		assert.ErrorIs(t, err, models.ErrOrderNotOpen)
	})

	t.Run("succesfully removes order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		order := models.Order{
			Id:     orderId,
			UserId: uuid.New(),
			Price:  price,
			Size:   size,
			Symbol: symbol,
			Side:   models.BUY,
			Status: models.STATUS_OPEN,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		userOrdersKey := CreateUserOrdersKey(order.UserId)
		orderIDKey := CreateOrderIDKey(order.Id)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetVal(1)
		mock.ExpectSRem(userOrdersKey, order.Id.String()).SetVal(1)
		mock.ExpectHSet(orderIDKey, "status", models.STATUS_CANCELED.String()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.RemoveOrder(ctx, order)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
