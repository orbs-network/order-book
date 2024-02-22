package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo_StoreNewPendingSwap(t *testing.T) {
	ctx := context.Background()

	// swapJson := []string{
	// 	`{"created": "2024-02-20T15:36:57.316643+02:00","frags":[{"orderId": "550e8400-e29b-41d4-a716-446655440000","inSize": "10.5"},{"orderId": "550e8400-e29b-41d4-a716-446655440001","inSize": "5.3"}]}`, `{"started":"*","created": "2024-02-20T15:36:57.316643+02:00","frags":[{"orderId": "550e8400-e29b-41d4-a716-446655440002","inSize": "7.8"}]}`,
	// }

	// swap := models.Swap{
	// 	Created:   time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
	// 	Frags:     []models.OrderFrag{},
	// 	Succeeded: false,
	// 	TxHash:    "0x123",
	// }

	// swapJson, _ := json.Marshal(swap)

	// t.Run("success - should add new pending swap to list given an existing swap is found and pending swap is not already tracked", func(t *testing.T) {
	// 	db, mock := redismock.NewClientMock()

	// 	repo := &redisRepository{
	// 		client: db,
	// 	}

	// 	pendingJson, _ := mocks.SwapTx.ToJson()
	// 	// var swap0 models.Swap
	// 	// fmt.Printf("swap0: %#v\n", swap0)
	// 	// fmt.Println("swap created", swap0.Created)
	// 	// err := json.Unmarshal([]byte(swapJson[0]), &swap0)
	// 	// assert.NoError(t, err)
	// 	//mock.ExpectSAdd(CreateSwapStartedKey(), mocks.SwapTx.SwapId.String()).SetVal(1)
	// 	//mock.ExpectSet(CreateSwapKey(mocks.SwapTx.SwapId), swapJson[0], 0) //.SetVal("1")
	// 	//mock.ExpectSet(CreatePendingSwapTxsKey(), swapJson[0], 0).SetVal("1")
	// 	///mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal(swapJson[0])
	// 	//mock.ExpectSet(CreateSwapKey(mocks.SwapTx.SwapId), swap0, 0).SetVal("1")
	// 	//mock.Regexp().ExpectSet(CreateSwapKey(mocks.SwapTx.SwapId), `\{(?:[^{}]|(?R))*\}`, 0) //.SetErr(errors.New("FAIL"))

	// 	mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal(string(swapJson))
	// 	//mock.ExpectSet(CreateSwapKey(mocks.SwapTx.SwapId), `^[0-9]+$`, 0).SetVal("1")
	// 	mock.ExpectSetArgs(CreateSwapKey(mocks.SwapTx.SwapId), swap, redis.SetArgs{Get: false}).SetVal("1")

	// 	//mock.Regexp().ExpectSet(CreateSwapKey(mocks.SwapTx.SwapId), `^[0-9]+$`, 0).SetVal("1")
	// 	//mock.Expect("SET").WithAnyArgs().WillReturn("OK")

	// 	// db.HSet(ctx, "key", "field", time.Now().Unix())
	// 	// mock.Regexp().ExpectHSet("key", "field", `^[0-9]+$`).SetVal(1)

	// 	mock.ExpectRPush(CreatePendingSwapTxsKey(), pendingJson).SetVal(1)

	// 	err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)
	// 	assert.NoError(t, err)
	// })

	t.Run("should return `ErrNotFound` error if existing swap not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal("")

		err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("should return error if redis fails to get swap", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetErr(assert.AnError)

		err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

		assert.ErrorContains(t, err, "failed to get swap unexpectedly")
	})

	// t.Run("should return `ErrValAlreadyInSet` error if swap already in tracker", func(t *testing.T) {
	// 	db, mock := redismock.NewClientMock()

	// 	repo := &redisRepository{
	// 		client: db,
	// 	}

	// 	mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal(swapJson[0])
	// 	//mock.ExpectLRange(CreateSwapKey(mocks.SwapTx.SwapId), 0, -1).SetVal(swapJson)
	// 	//mock.ExpectSAdd(CreateSwapStartedKey(), mocks.SwapTx.SwapId.String()).SetErr(models.ErrValAlreadyInSet)

	// 	err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

	// 	assert.ErrorIs(t, err, models.ErrValAlreadyInSet)
	// })

	// t.Run("should return error if redis fails to add pending swap to tracker", func(t *testing.T) {
	// 	db, mock := redismock.NewClientMock()

	// 	repo := &redisRepository{
	// 		client: db,
	// 	}

	// 	mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal(swapJson[0])
	// 	//mock.ExpectLRange(CreateSwapKey(mocks.SwapTx.SwapId), 0, -1).SetVal(swapJson)
	// 	//mock.ExpectSAdd(CreateSwapStartedKey(), mocks.SwapTx.SwapId.String()).SetErr(assert.AnError)

	// 	err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

	// 	assert.ErrorContains(t, err, "failed to add pendingSwap to tracker")
	// })

	// t.Run("should return error if redis fails to store pending swap", func(t *testing.T) {
	// 	db, mock := redismock.NewClientMock()

	// 	repo := &redisRepository{
	// 		client: db,
	// 	}

	// 	mock.ExpectGet(CreateSwapKey(mocks.SwapTx.SwapId)).SetVal(swapJson[0])
	// 	//mock.ExpectLRange(CreateSwapKey(mocks.SwapTx.SwapId), 0, -1).SetVal(swapJson)
	// 	//mock.ExpectSAdd(CreateSwapStartedKey(), mocks.SwapTx.SwapId.String()).SetVal(1)
	// 	mock.ExpectRPush(CreatePendingSwapTxsKey(), mocks.SwapTx.ToMap()).SetErr(assert.AnError)

	// 	err := repo.StoreNewPendingSwap(ctx, mocks.SwapTx)

	// 	assert.ErrorContains(t, err, "failed to store pending swap tx")
	// })
}
