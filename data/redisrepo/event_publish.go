package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// PublishEvent publishes an event to Redis
func (r *redisRepository) PublishEvent(ctx context.Context, key string, value interface{}) error {
	err := r.client.Publish(ctx, key, value).Err()

	if err != nil {
		logctx.Error(ctx, "failed to publish redis event", logger.Error(err), logger.String("key", key))
		return fmt.Errorf("failed to publish redis event: %v", err)
	}

	logctx.Info(ctx, "published redis event", logger.String("key", key))
	return nil
}
