package redisrepo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_StoreCompletedSwap(t *testing.T) {

	swapId := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	txId := "0xsomething"
	timestamp := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

	mockCompletedSwap := store.StoreCompletedSwapInput{
		UserId:    userId,
		SwapId:    swapId,
		OrderId:   orderId,
		TxId:      txId,
		Timestamp: timestamp,
		Block:     123,
	}

	json, _ := json.Marshal(mockCompletedSwap)

	t.Run("store completed swap success", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreateCompletedSwapsKey(mockCompletedSwap.UserId), json).SetVal(1)

		err := repo.StoreCompletedSwap(ctx, mockCompletedSwap)

		assert.NoError(t, err, "should not return error")
	})

	t.Run("store completed swap failure", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectRPush(CreateCompletedSwapsKey(mockCompletedSwap.UserId), json).SetErr(assert.AnError)

		err := repo.StoreCompletedSwap(ctx, mockCompletedSwap)

		assert.ErrorContains(t, err, "failed to store completed swap in Redis")
	})

}
