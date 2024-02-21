package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreCompletedSwap(ctx context.Context, input store.StoreCompletedSwapInput) error {
	key := CreateCompletedSwapsKey(input.UserId)

	completedSwapJson, err := json.Marshal(input)
	if err != nil {
		logctx.Error(ctx, "failed to marshal completed swap", logger.Error(err), logger.String("userId", input.UserId.String()), logger.String("swapId", input.SwapId.String()), logger.String("txHash", input.TxId))
		return fmt.Errorf("failed to marshal completed swap: %v", err)
	}

	_, err = r.client.RPush(ctx, key, completedSwapJson).Result()
	if err != nil {
		logctx.Error(ctx, "failed to store completed swap in Redis", logger.Error(err), logger.String("userId", input.UserId.String()), logger.String("swapId", input.SwapId.String()), logger.String("txHash", input.TxId))
		return fmt.Errorf("failed to store completed swap in Redis: %v", err)
	}

	return nil
}
