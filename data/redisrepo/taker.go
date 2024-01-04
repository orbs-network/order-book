package redisrepo

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) GetSwap(ctx context.Context, swapId uuid.UUID) ([]models.OrderFrag, error) {
	swapKey := CreateSwapKey(swapId)

	swapJsons, err := r.client.LRange(ctx, "uvix"+swapKey, 0, -1).Result()
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
