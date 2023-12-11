package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreAuction(ctx context.Context, auctionID uuid.UUID, frags []models.OrderFrag) error {

	auctionJson, err := models.MarshalOrderFrags(frags)
	if err != nil {
		logctx.Error(ctx, "failed to marshal auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
		return fmt.Errorf("failed to marshal auction: %v", err)
	}

	auctionKey := CreateAuctionKey(auctionID)

	_, err = r.client.RPush(ctx, auctionKey, auctionJson).Result()
	if err != nil {
		logctx.Error(ctx, "failed to store auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
		return fmt.Errorf("failed to store auction: %v", err)
	}

	// Set the TTL to 24 hours (24 hours * 60 minutes * 60 seconds)
	// TODO:
	// err = r.client.Expire(ctx, auctionKey, 24*time.Hour).Err()
	// if err != nil {
	// 	fmt.Println("Error setting key:", err)
	// 	return models.ErrUnexpectedError
	// }

	logctx.Info(ctx, "stored auction", logger.String("auctionID", auctionID.String()))
	return nil
}

func (r *redisRepository) StoreSwap(ctx context.Context, swapId uuid.UUID, frags []models.OrderFrag) error {

	swapJson, err := models.MarshalOrderFrags(frags)
	if err != nil {
		logctx.Error(ctx, "failed to marshal swap", logger.String("swapID", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to marshal swap: %v", err)
	}

	auctionKey := CreateSwapKey(swapId)

	_, err = r.client.RPush(ctx, auctionKey, swapJson).Result()
	if err != nil {
		logctx.Error(ctx, "failed to store auction", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to store auction: %v", err)
	}

	// Set the TTL to 24 hours (24 hours * 60 minutes * 60 seconds)
	// TODO:
	// err = r.client.Expire(ctx, auctionKey, 24*time.Hour).Err()
	// if err != nil {
	// 	fmt.Println("Error setting key:", err)
	// 	return models.ErrUnexpectedError
	// }

	logctx.Info(ctx, "stored auction", logger.String("swapId", swapId.String()))
	return nil
}

func (r *redisRepository) RemoveSwap(ctx context.Context, swapId uuid.UUID) error {
	auctionKey := CreateSwapKey(swapId)
	err := r.client.Del(ctx, auctionKey).Err()
	logctx.Error(ctx, "Redis del failed", logger.String("key", auctionKey), logger.Error(err))
	return err
}
