package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) SubscribeToEvents(ctx context.Context, channel string) (chan []byte, error) {
	logctx.Info(ctx, "subscribing to channel", logger.String("channel", channel))

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

	// Start a goroutine to listen for messages
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for msg := range ch {
			messages <- []byte(msg.Payload)
		}
		logctx.Info(ctx, "subscription ended", logger.String("channel", channel))
		close(messages) // Close the channel when the subscription ends
	}()

	return messages, nil
}
