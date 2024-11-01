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

func (r *redisRepository) saveSwap(ctx context.Context, swapId uuid.UUID, swap models.Swap, resolved bool) error {
	swapJson, err := json.Marshal(swap)
	if err != nil {
		logctx.Error(ctx, "failed to marshal swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to marshal swap: %v", err)
	}

	swapKey := CreateOpenSwapKey(swapId)
	if resolved {
		swapKey = CreateResolvedSwapKey(swapId)
	}

	_, err = r.client.Set(ctx, swapKey, swapJson, 0).Result()

	return err
}

func (r *redisRepository) StoreSwap(ctx context.Context, swapId uuid.UUID, symbol models.Symbol, side models.Side, frags []models.OrderFrag) error {
	swap := models.NewSwap(symbol, side, frags)

	err := r.saveSwap(ctx, swapId, *swap, false)
	if err != nil {
		logctx.Error(ctx, "failed to store swap", logger.String("swapId", swapId.String()), logger.Error(err))
		return fmt.Errorf("failed to store swap: %v", err)
	}
	logctx.Info(ctx, "stored swap", logger.String("swapId", swapId.String()))
	return nil
}

func (r *redisRepository) GetSwap(ctx context.Context, swapId uuid.UUID, open bool) (*models.Swap, error) {
	swapKey := CreateOpenSwapKey(swapId)
	if !open {
		swapKey = CreateResolvedSwapKey(swapId)
	}

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

	return &swap, nil
}

func (r *redisRepository) RemoveSwap(ctx context.Context, swapId uuid.UUID) error {
	logctx.Debug(ctx, "RemoveSwap", logger.String("key", swapId.String()))
	swapKey := CreateOpenSwapKey(swapId)
	err := r.client.Del(ctx, swapKey).Err()
	if err != nil {
		logctx.Error(ctx, "RemoveSwap Redis del failed", logger.String("key", swapKey), logger.Error(err))
	}
	return err
}
