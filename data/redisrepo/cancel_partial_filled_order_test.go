package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_CancelPartialFilledOrder(t *testing.T) {

	partialFilledOrder := models.Order{
		Id:         orderId,
		Price:      price,
		Size:       size,
		SizeFilled: size.Div(decimal.NewFromFloat(2)),
		Symbol:     symbol,
		Side:       models.BUY,
	}

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

		err := repo.CancelPartialFilledOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrOrderPending)
	})

	t.Run("rejects unfilled order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		order := models.Order{
			Id:     orderId,
			Price:  price,
			Size:   size,
			Symbol: symbol,
			Side:   models.BUY,
		}

		mock.ExpectTxPipeline()

		err := repo.CancelPartialFilledOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrOrderNotPartialFilled)
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
			SizeFilled:  size,
			SizePending: decimal.Zero,
			Symbol:      symbol,
			Side:        models.BUY,
		}

		mock.ExpectTxPipeline()

		err := repo.CancelPartialFilledOrder(ctx, order)

		assert.ErrorIs(t, err, models.ErrOrderNotPartialFilled)
	})

	t.Run("relevant error returned when fails to cancel order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(partialFilledOrder.Id), "cancelled", "true").SetErr(assert.AnError)

		err := repo.CancelPartialFilledOrder(ctx, partialFilledOrder)

		assert.ErrorContains(t, err, "failed to cancel partial filled order")
	})

	t.Run("succesfully cancels order partial order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(CreateOrderIDKey(partialFilledOrder.Id), "cancelled", "true").SetVal(1)
		mock.ExpectZRem(CreateBuySidePricesKey(partialFilledOrder.Symbol), partialFilledOrder.Id.String()).SetVal(1)
		mock.ExpectZRem(CreateUserOpenOrdersKey(partialFilledOrder.UserId), partialFilledOrder.Id.String()).SetVal(1)
		mock.ExpectTxPipelineExec()

		err := repo.CancelPartialFilledOrder(ctx, partialFilledOrder)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
