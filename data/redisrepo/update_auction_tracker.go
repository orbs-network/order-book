package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) UpdateAuctionTracker(ctx context.Context, auctionStatus models.AuctionStatus, auctionId uuid.UUID) error {

	auctionTrackerKey := CreateAuctionTrackerKey(auctionStatus)

	if err := AddVal2Set(ctx, r.client, auctionTrackerKey, auctionId.String()); err != nil {
		if err == models.ErrValAlreadyInSet {
			logctx.Warn(ctx, "auction already in tracker", logger.String("auctionId", auctionId.String()), logger.String("auctionStatus", auctionStatus.String()))
			return err
		}

		logctx.Error(ctx, "failed to add auction to tracker", logger.Error(err), logger.String("auctionId", auctionId.String()), logger.String("auctionStatus", auctionStatus.String()))
		return fmt.Errorf("failed to add auction to tracker: %w", err)
	}

	logctx.Info(ctx, "added auction to tracker", logger.String("auctionId", auctionId.String()), logger.String("auctionStatus", auctionStatus.String()))
	return nil
}
