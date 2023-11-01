package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) AddVal2Set(ctx context.Context, key, val string) error {
	isMember, err := r.client.SIsMember(ctx, key, val).Result()
	if err != nil {
		logctx.Warn(ctx, "SIsMember Failed", logger.Error(err))
		return err
	}
	if isMember {
		logctx.Warn(ctx, err.Error())
		return models.ErrValAlreadyInSet
	}

	_, err = r.client.SAdd(ctx, key, val).Result()
	if err != nil {
		logctx.Warn(ctx, "SAdd Failed", logger.Error(err))
		return err
	}

	return nil

}
