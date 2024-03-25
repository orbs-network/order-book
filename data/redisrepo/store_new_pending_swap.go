package redisrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// StoreNewPendingSwap stores a new pending swap in order for its status (pending/complete) to be checked later
func (r *redisRepository) StoreNewPendingSwap(ctx context.Context, p models.SwapTx) error {
	// confirm swapID is valid
	swap, err := r.GetSwap(ctx, p.SwapId, true)
	if err != nil {
		if err == models.ErrNotFound {
			logctx.Warn(ctx, "no swap found by that ID", logger.Error(err), logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
			return err
		}
		logctx.Error(ctx, "failed to get swap", logger.Error(err), logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
		return fmt.Errorf("failed to get swap unexpectedly: %s", err)
	}

	// protect re-entry
	if swap.IsStarted() {
		logctx.Error(ctx, "swap is already started", logger.Error(err), logger.String("startedSwapId", p.SwapId.String()))
		return fmt.Errorf("swap is already started: %s", err)
	}

	// update swapId:started field
	swap.Started = time.Now()
	swap.TxHash = p.TxHash
	swap.Id = p.SwapId
	err = r.saveSwap(ctx, p.SwapId, *swap, false)
	if err != nil {
		logctx.Error(ctx, "failed to update swap started time", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return fmt.Errorf("failed to update swap started time: %s", err)
	}

	logctx.Debug(ctx, "store pending swap", logger.String("swapId", p.SwapId.String()), logger.String("txHash", p.TxHash))
	return nil
}
