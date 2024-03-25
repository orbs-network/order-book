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

	logctx.Debug(ctx, "published redis event", logger.String("key", key))
	return nil
}

// SubscribeToEvents subscribes to events on a given Redis channel
func (r *redisRepository) SubscribeToEvents(ctx context.Context, channel string) (chan []byte, error) {
	logctx.Debug(ctx, "subscribing to channel", logger.String("channel", channel))

	// Subscribe to the specified channel
	pubsub := r.client.Subscribe(ctx, channel)

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		logctx.Error(ctx, "error on receiving from pubsub", logger.Error(err), logger.String("channel", channel))
		return nil, fmt.Errorf("error on receiving from pubsub: %w", err)
	}

	// Create a channel to pass messages to the caller
	messages := make(chan []byte)

	// Listen for messages
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for msg := range ch {
			messages <- []byte(msg.Payload)
		}
		logctx.Debug(ctx, "subscription ended", logger.String("channel", channel))
		close(messages)
	}()

	return messages, nil
}
