package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreSwap(ctx context.Context, swapId uuid.UUID, frags []models.OrderFrag) error {

	swapJson, err := models.MarshalOrderFrags(frags)
	if err != nil {
		logctx.Error(ctx, "failed to marshal swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to marshal swap: %v", err)
	}

	swapKey := CreateSwapKey(swapId)

	_, err = r.client.RPush(ctx, swapKey, swapJson).Result()
	if err != nil {
		logctx.Error(ctx, "failed to store swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to store swap: %v", err)
	}

	// Set the TTL to 24 hours (24 hours * 60 minutes * 60 seconds)
	// TODO:
	// err = r.client.Expire(ctx, swapKey, 24*time.Hour).Err()
	// if err != nil {
	// 	fmt.Println("Error setting key:", err)
	// 	return models.ErrUnexpectedError
	// }

	logctx.Info(ctx, "stored swap", logger.String("swapId", swapId.String()))
	return nil
}
