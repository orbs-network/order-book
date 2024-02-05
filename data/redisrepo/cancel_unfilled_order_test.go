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
var symbol, _ = models.StrToSymbol("MATIC-USDC")
var price = decimal.NewFromFloat(10.0)

var order = models.Order{
	Id:     orderId,
	Price:  price,
	Size:   size,
	Symbol: symbol,
	Side:   models.BUY,
}

func TestRedisRepository_CancelUnfilledOrder(t *testing.T) {

	t.Run("rejects pending order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		order := models.Order{
			Id:          orderId,
			Price:       price,
			Size:        size,
			SizePending: size.Div(decimal.NewFromFloat(2)),
			Symbol:      symbol,
			Side:        models.BUY,
		}

		mock.ExpectTxPipeline()

		err := repo.CancelUnfilledOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrOrderPending)
	})

	t.Run("rejects filled order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		order := models.Order{
			Id:          orderId,
			Price:       price,
			Size:        size,
			SizePending: decimal.Zero,
			SizeFilled:  size,
			Symbol:      symbol,
			Side:        models.BUY,
		}

		mock.ExpectTxPipeline()

		err := repo.CancelUnfilledOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrOrderNotUnfilled)
	})

	t.Run("fails to remove order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetErr(assert.AnError)

		err := repo.CancelUnfilledOrder(ctx, order)

		assert.ErrorContains(t, err, "failed to remove order")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("succesfully removes order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		buyPricesKey := CreateBuySidePricesKey(order.Symbol)
		userOrdersKey := CreateUserOpenOrdersKey(order.UserId)
		clientOIdKey := CreateClientOIDKey(order.ClientOId)
		orderIdKey := CreateOrderIDKey(order.Id)

		mock.ExpectTxPipeline()
		mock.ExpectZRem(buyPricesKey, order.Id.String()).SetVal(1)
		mock.ExpectZRem(userOrdersKey, order.Id.String()).SetVal(1)
		mock.ExpectDel(clientOIdKey).SetVal(1)
		mock.ExpectDel(orderIdKey).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.CancelUnfilledOrder(ctx, order)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
