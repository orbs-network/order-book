package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// ProcessCompletedSwapOrders stores the updated swap orders and removes the swap from Redis. It should be called after a swap is completed.
//
// `orders` should be the orders that were part of the swap (with `SizePending` and `SizeFilled` updated accordingly)
//
// `isSuccessful` should be `true` if the swap was successful, `false` otherwise

func (r *redisRepository) ResolveSwap(ctx context.Context, swap models.Swap) error {

	// save swap in resolved key
	err := r.saveSwap(ctx, swap.Id, swap, true)
	if err != nil {
		logctx.Error(ctx, "failed to save swap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return err
	}

	// remove from swapId
	err = r.RemoveSwap(ctx, swap.Id)
	if err != nil {
		logctx.Error(ctx, "failed to remove swap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return err
	}

	return nil
}

// save swapId in a set of the userId:resolvedSwap key
func (r *redisRepository) StoreUserResolvedSwap(ctx context.Context, userId uuid.UUID, swap models.Swap) error {
	key := CreateUserResolvedSwapsKey(userId)
	return AddVal2Set(ctx, r.client, key, swap.Id.String())
}

func (r *redisRepository) GetUserResolvedSwapIds(ctx context.Context, userId uuid.UUID) ([]string, error) {
	key := CreateUserResolvedSwapsKey(userId)
	res, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, models.ErrNotFound
		}
		logctx.Error(ctx, "could not get user swaps", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, err
	}
	return res, nil
}
