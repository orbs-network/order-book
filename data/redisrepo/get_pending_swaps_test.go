package redisrepo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func generatePendingSwapStr(num int) []string {
	var pendingSwaps []string
	for i := 0; i < num; i++ {
		pendingSwaps = append(pendingSwaps, "{\"swapId\":\"75fff27c-5a98-4f04-9bba-ce932d9c8e67\",\"txHash\":\"cc6fb26e-bc79-4ea3-b761-224893dc7df3\"}")
	}
	return pendingSwaps
}

func TestRedisRepo_GetPendingSwaps(t *testing.T) {
	ctx := context.Background()

	t.Run("should get a list of pending swaps", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectLRange(CreatePendingSwapTxsKey(), 0, -1).SetVal([]string{
			"{\"swapId\":\"75fff27c-5a98-4f04-9bba-ce932d9c8e67\",\"txHash\":\"0xefd319bb86b954a8e8cd7d9396546db8d3251910209cd8b1b9a674ef8585f226\"}",
		})

		pendingSwaps, err := repo.GetPendingSwaps(ctx)

		assert.NoError(t, err)
		assert.Len(t, pendingSwaps, 1)
		assert.Equal(t, uuid.MustParse("75fff27c-5a98-4f04-9bba-ce932d9c8e67"), pendingSwaps[0].SwapId)
		assert.Equal(t, "0xefd319bb86b954a8e8cd7d9396546db8d3251910209cd8b1b9a674ef8585f226", pendingSwaps[0].TxHash)
	})

	t.Run("should work with a large list of pending swaps", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		numSwaps := 10000

		mock.ExpectLRange(CreatePendingSwapTxsKey(), 0, -1).SetVal(generatePendingSwapStr(numSwaps))

		pendingSwaps, err := repo.GetPendingSwaps(ctx)

		assert.NoError(t, err)
		assert.Len(t, pendingSwaps, numSwaps)
	})

	t.Run("should return an empty list if there are no pending swaps", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectLRange(CreatePendingSwapTxsKey(), 0, -1).SetVal([]string{})

		pendingSwaps, err := repo.GetPendingSwaps(ctx)

		assert.NoError(t, err)
		assert.Len(t, pendingSwaps, 0)
	})

	t.Run("should return an error if there is an error getting the pending swaps", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectLRange(CreatePendingSwapTxsKey(), 0, -1).SetErr(assert.AnError)

		_, err := repo.GetPendingSwaps(ctx)

		assert.ErrorContains(t, err, "failed to get pending swaps")
	})

}
