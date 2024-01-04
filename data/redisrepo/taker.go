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

func (r *redisRepository) GetSwap(ctx context.Context, swapId uuid.UUID) ([]models.OrderFrag, error) {
	swapKey := CreateSwapKey(swapId)

	swapJsons, err := r.client.LRange(ctx, swapKey, 0, -1).Result()
	// Error
	if err != nil {
		logctx.Error(ctx, "failed to get swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return []models.OrderFrag{}, models.ErrUnexpectedError
	}
	// empty range of swaps
	if len(swapJsons) == 0 {
		logctx.Error(ctx, "swap is not found", logger.String("swapId", swapId.String()), logger.Error(err))
		return []models.OrderFrag{}, models.ErrNotFound
	}

	var frags []models.OrderFrag
	for _, swapJson := range swapJsons {
		var orders []models.OrderFrag
		err := json.Unmarshal([]byte(swapJson), &orders)
		if err != nil {
			logctx.Error(ctx, "failed to unmarshal swap", logger.String("swapId", swapId.String()), logger.Error(err))
			return []models.OrderFrag{}, models.ErrMarshalError
		}
		frags = append(frags, orders...)
	}

	logctx.Info(ctx, "got swap", logger.String("swapId", swapId.String()))
	return frags, nil
}

func (r *redisRepository) RemoveSwap(ctx context.Context, swapId uuid.UUID) error {
	swapKey := CreateSwapKey(swapId)
	err := r.client.Del(ctx, swapKey).Err()
	logctx.Error(ctx, "Redis del failed", logger.String("key", swapKey), logger.Error(err))
	return err
}
