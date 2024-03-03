package redisrepo

import (
	"testing"
)

func TestRedisRepo_ResolveSwap(t *testing.T) {
	//ctx := context.Background()

	// db, mock := redismock.NewClientMock()

	// repo := &redisRepository{
	// 	client: db,
	// }

	// swapId := uuid.MustParse("00000000-0000-0000-0000-000000000007")
	// timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	// block := int64(1234567)

	// mockTx := models.Tx{
	// 	Status:    models.TX_SUCCESS,
	// 	TxHash:    "0x123",
	// 	Block:     &block,
	// 	Timestamp: &timestamp,
	// }

	t.Run("should process successful swap with completely filled order and partially filled order", func(t *testing.T) {

		// filledOrder := models.Order{
		// 	SizeFilled: decimal.NewFromFloat(100),
		// 	Size:       decimal.NewFromFloat(100),
		// 	Side:       models.SELL,
		// 	Timestamp:  timestamp,
		// }

		// partiallyFilledOrder := models.Order{
		// 	SizeFilled: decimal.NewFromFloat(50),
		// 	Size:       decimal.NewFromFloat(100),
		// 	Side:       models.SELL,
		// 	Timestamp:  timestamp,
		// }

		// mock.ExpectTxPipeline()

		// // Store completely filled order
		// mock.ExpectZRem(CreateUserOpenOrdersKey(filledOrder.UserId), filledOrder.Id.String()).SetVal(1)
		// mock.ExpectZAdd(CreateUserFilledOrdersKey(filledOrder.UserId), redis.Z{
		// 	Score:  float64(timestamp.UTC().UnixNano()),
		// 	Member: filledOrder.Id.String(),
		// }).SetVal(1)
		// mock.ExpectZRem(CreateSellSidePricesKey(filledOrder.Symbol), filledOrder.Id.String()).SetVal(1)
		// mock.ExpectHSet(CreateOrderIDKey(filledOrder.Id), filledOrder.OrderToMap()).SetVal(1)

		// // Store partially filled order
		// mock.ExpectZAdd(CreateUserOpenOrdersKey(partiallyFilledOrder.UserId), redis.Z{
		// 	Score:  float64(timestamp.UTC().UnixNano()),
		// 	Member: partiallyFilledOrder.Id.String(),
		// }).SetVal(1)
		// mock.ExpectHSet(CreateOrderIDKey(partiallyFilledOrder.Id), partiallyFilledOrder.OrderToMap()).SetVal(1)
		// mock.ExpectSet(CreateClientOIDKey(partiallyFilledOrder.ClientOId), partiallyFilledOrder.Id.String(), 0).SetVal("OK")
		// f64Price, _ := partiallyFilledOrder.Price.Float64()
		// timestamp := float64(partiallyFilledOrder.Timestamp.UTC().UnixNano()) / 1e9
		// score := f64Price + (timestamp / 1e12)
		// mock.ExpectZAdd(CreateSellSidePricesKey(partiallyFilledOrder.Symbol), redis.Z{
		// 	Score:  score,
		// 	Member: partiallyFilledOrder.Id.String(),
		// }).SetVal(1)

		// // Remove swap
		// mock.ExpectDel(CreateOpenSwapKey(swapId)).SetVal(1)

		// mock.ExpectTxPipelineExec()

		// err := repo.RemoveSwap(ctx, swapId)
		// assert.NoError(t, err)
	})

	t.Run("should process failed swap", func(t *testing.T) {
		// orderToBeRolledback := models.Order{
		// 	SizeFilled:  decimal.NewFromFloat(20),
		// 	Size:        decimal.NewFromFloat(100),
		// 	SizePending: decimal.NewFromFloat(0),
		// 	Side:        models.SELL,
		// 	Timestamp:   timestamp,
		// }

		// mock.ExpectTxPipeline()

		// // Store order
		// mock.ExpectZAdd(CreateUserOpenOrdersKey(orderToBeRolledback.UserId), redis.Z{
		// 	Score:  float64(timestamp.UTC().UnixNano()),
		// 	Member: orderToBeRolledback.Id.String(),
		// }).SetVal(1)
		// mock.ExpectHSet(CreateOrderIDKey(orderToBeRolledback.Id), orderToBeRolledback.OrderToMap()).SetVal(1)
		// mock.ExpectSet(CreateClientOIDKey(orderToBeRolledback.ClientOId), orderToBeRolledback.Id.String(), 0).SetVal("OK")
		// f64Price, _ := orderToBeRolledback.Price.Float64()
		// timestamp := float64(orderToBeRolledback.Timestamp.UTC().UnixNano()) / 1e9
		// score := f64Price + (timestamp / 1e12)
		// mock.ExpectZAdd(CreateSellSidePricesKey(orderToBeRolledback.Symbol), redis.Z{
		// 	Score:  score,
		// 	Member: orderToBeRolledback.Id.String(),
		// }).SetVal(1)

		// // Remove swap
		// mock.ExpectDel(CreateOpenSwapKey(swapId)).SetVal(1)

		// mock.ExpectTxPipelineExec()

		// err := repo.ProcessCompletedSwapOrders(ctx, []store.OrderWithSize{{
		// 	Order: &orderToBeRolledback,
		// 	Size:  decimal.NewFromFloat(20),
		// }}, swapId, &mockTx, false)
		// assert.NoError(t, err)
	})

	t.Run("no database writes should happen if part of transaction fails", func(t *testing.T) {

		// mock.ExpectTxPipeline()
		// mock.ExpectZAdd(CreateUserOpenOrdersKey(mocks.Order.UserId), redis.Z{
		// 	Score:  float64(timestamp.UTC().UnixNano()),
		// 	Member: mocks.Order.Id.String(),
		// }).SetErr(assert.AnError)

		// err := repo.ProcessCompletedSwapOrders(ctx, []store.OrderWithSize{
		// 	{
		// 		Order: &mocks.Order,
		// 		Size:  decimal.NewFromFloat(20),
		// 	},
		// }, swapId, &mockTx, true)
		// assert.ErrorContains(t, err, "failed to execute ProcessCompletedSwapOrders transaction")
	})

}
