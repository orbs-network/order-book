package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StorePendingSwaps(ctx context.Context, pendingSwaps []models.Pending) error {
	key := CreatePendingSwapTxsKey()

	transaction := r.client.TxPipeline()

	// Delete old pending swaps list
	if _, err := transaction.Del(ctx, key).Result(); err != nil {
		logctx.Error(ctx, "failed to delete old pending swaps", logger.Error(err))
		return fmt.Errorf("failed to delete old pending swaps: %s", err)
	}

	// Only queue RPUSH command if there are pending swaps
	if len(pendingSwaps) > 0 {
		pendingSwapsJson := make([]interface{}, 0, len(pendingSwaps))
		for _, p := range pendingSwaps {
			jsonData, err := json.Marshal(p) // Ensure this marshals to JSON properly
			if err != nil {
				logctx.Error(ctx, "failed to marshal pending swap", logger.Error(err))
				return fmt.Errorf("failed to marshal pending swap: %s", err)
			}

			pendingSwapsJson = append(pendingSwapsJson, string(jsonData))
		}

		if _, err := transaction.RPush(ctx, key, pendingSwapsJson...).Result(); err != nil {
			logctx.Error(ctx, "failed to push pending swaps to Redis", logger.Error(err))
			return fmt.Errorf("failed to push pending swaps to Redis: %s", err)
		}
	}

	// Execute the transaction
	_, err := transaction.Exec(ctx)
	if err != nil {
		logctx.Error(ctx, "failed to execute Redis transaction", logger.Error(err))
		return fmt.Errorf("failed to execute Redis transaction: %s", err)
	}

	return nil
}
