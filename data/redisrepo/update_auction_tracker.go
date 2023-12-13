package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) UpdateSwapTracker(ctx context.Context, swapStatus models.SwapStatus, swapId uuid.UUID) error {

	swapTrackerKey := CreateSwapTrackerKey(swapStatus)

	if err := AddVal2Set(ctx, r.client, swapTrackerKey, swapId.String()); err != nil {
		if err == models.ErrValAlreadyInSet {
			logctx.Warn(ctx, "swap already in tracker", logger.String("swapId", swapId.String()), logger.String("swapStatus", swapStatus.String()))
			return err
		}

		logctx.Error(ctx, "failed to add swap to tracker", logger.Error(err), logger.String("swapId", swapId.String()), logger.String("swapStatus", swapStatus.String()))
		return fmt.Errorf("failed to add swap to tracker: %w", err)
	}

	logctx.Info(ctx, "added swap to tracker", logger.String("swapId", swapId.String()), logger.String("swapStatus", swapStatus.String()))
	return nil
}
