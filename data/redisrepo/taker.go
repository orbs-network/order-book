package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreSwap(ctx context.Context, swapId uuid.UUID, frags []models.OrderFrag) error {
	swap := models.NewSwap(frags)
	swapJson, err := json.Marshal(swap)
	if err != nil {
		logctx.Error(ctx, "failed to marshal swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to marshal swap: %v", err)
	}

	swapKey := CreateSwapKey(swapId)

	//_, err = r.client.RPush(ctx, swapKey, swapJson).Result()
	_, err = r.client.Set(ctx, swapKey, swapJson, 0).Result()

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

func (r *redisRepository) GetSwap(ctx context.Context, swapId uuid.UUID) (*models.Swap, error) {
	swapKey := CreateSwapKey(swapId)

	swapJson, err := r.client.Get(ctx, swapKey).Result()
	// Error
	if err != nil {
		logctx.Error(ctx, "failed to get swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return nil, models.ErrUnexpectedError
	}
	// empty range of swaps
	if len(swapJson) == 0 {
		logctx.Error(ctx, "swap is not found", logger.String("swapId", swapId.String()), logger.Error(err))
		return nil, models.ErrNotFound
	}

	var swap models.Swap
	err = json.Unmarshal([]byte(swapJson), &swap)
	if err != nil {
		logctx.Error(ctx, "failed to unmarshal swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return nil, models.ErrMarshalError
	}

	logctx.Info(ctx, "got swap", logger.String("swapId", swapId.String()))
	return &swap, nil
}

func (r *redisRepository) RemoveSwap(ctx context.Context, swapId uuid.UUID) error {
	swapKey := CreateSwapKey(swapId)
	err := r.client.Del(ctx, swapKey).Err()
	logctx.Error(ctx, "Redis del failed", logger.String("key", swapKey), logger.Error(err))
	return err
}
