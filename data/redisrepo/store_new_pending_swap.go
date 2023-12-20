package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// StoreNewPendingSwap stores a new pending swap in order for its status (pending/complete) to be checked later
func (r *redisRepository) StoreNewPendingSwap(ctx context.Context, p models.SwapTx) error {
	key := CreatePendingSwapTxsKey()

	jsonData, err := json.Marshal(p)
	if err != nil {
		logctx.Error(ctx, "failed to marshal pending", logger.Error(err), logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
		return fmt.Errorf("failed to marshal pending: %s", err)
	}

	_, err = r.client.RPush(ctx, key, jsonData).Result()

	if err != nil {
		logctx.Error(ctx, "failed to store pending swap tx", logger.Error(err), logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
		return fmt.Errorf("failed to store pending swap tx: %s", err)
	}

	logctx.Info(ctx, "stored pending swap tx", logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
	return nil
}
