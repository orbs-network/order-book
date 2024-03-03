package redisrepo

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func Key2Id(key string) string {
	splt := strings.Split(key, ":")

	if len(splt) < 2 {
		return ""
	}
	return splt[2]
}

func Key2UUID(ctx context.Context, key string) *uuid.UUID {
	id := Key2Id(key)
	if id == "" {
		logctx.Error(ctx, "swap key Invalid", logger.String("key", key))
		return nil
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		logctx.Error(ctx, "swap key Invalid", logger.Error(err), logger.String("key", key))
		return nil
	}
	return &uid
}

func (r *redisRepository) GetOpenSwaps(ctx context.Context) ([]models.Swap, error) {
	res := []models.Swap{}
	keys, err := r.EnumSubKeysOf(ctx, "swap:open")
	if err != nil {
		logctx.Error(ctx, "Failed to enum swap:open keys", logger.Error(err))
		return res, err
	}
	for _, key := range keys {
		id := Key2UUID(ctx, key)
		if id != nil {
			swap, err := r.GetSwap(ctx, *id)
			swap.Id = *id // confirm exists
			if err != nil {
				logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapKey", key))
			} else {
				// swap was started but not resolved
				if !swap.Started.IsZero() && swap.Resolved.IsZero() {
					// make sure id is there
					res = append(res, *swap)
				}
			}
		} else {
			logctx.Error(ctx, "failed to create id from key", logger.String("key", key))
		}
	}
	return res, nil
}
