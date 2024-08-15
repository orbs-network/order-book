package redisrepo

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

var BUFFER_SIZE = os.Getenv("REDIS_BUFFER_SIZE")

// PublishEvent publishes an event to a Redis channel
func (r *redisRepository) PublishEvent(ctx context.Context, key string, value interface{}) error {
	err := r.client.Publish(ctx, key, value).Err()

	if err != nil {
		logctx.Error(ctx, "failed to publish redis event", logger.Error(err), logger.String("key", key))
		return fmt.Errorf("failed to publish redis event: %v", err)
	}
	return nil
}

func (r *redisRepository) SubscribeToEvents(ctx context.Context, channel string) (chan []byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, exists := r.subscriptions[channel]
	if !exists {
		pubsub := r.client.Subscribe(ctx, channel)
		_, err := pubsub.Receive(ctx)
		if err != nil {
			return nil, fmt.Errorf("error subscribing to channel: %w", err)
		}

		sub = &channelSubscription{
			pubsub:   pubsub,
			clients:  make(map[chan []byte]struct{}),
			messages: make(chan []byte, getBufferSize()),
		}
		r.subscriptions[channel] = sub

		// Start a goroutine to distribute messages to clients
		go r.distributeMessages(ctx, channel, sub)
	}

	// Create a new client channel and add it to the subscription's clients
	clientChan := make(chan []byte, getBufferSize())
	sub.clients[clientChan] = struct{}{}

	return clientChan, nil
}

// UnsubscribeFromEvents unsubscribes a client from a Redis channel
func (r *redisRepository) UnsubscribeFromEvents(ctx context.Context, channel string, clientChan chan []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sub, exists := r.subscriptions[channel]
	if !exists {
		return
	}

	delete(sub.clients, clientChan)
	close(clientChan)

	if len(sub.clients) == 0 {
		if err := sub.pubsub.Close(); err != nil {
			logctx.Error(ctx, "error closing pubsub", logger.Error(err), logger.String("channel", channel))
		}
		delete(r.subscriptions, channel)
	}
}

// distributeMessages listens for messages on a Redis channel and sends them to all subscribed clients
func (r *redisRepository) distributeMessages(ctx context.Context, channel string, sub *channelSubscription) {
	defer func() {
		sub.pubsub.Close()
		close(sub.messages)
		r.mu.Lock()
		delete(r.subscriptions, channel)
		r.mu.Unlock()
	}()

	for {
		select {
		case msg := <-sub.pubsub.Channel():
			// Send the message to all subscribed clients
			r.mu.Lock()

			if len(sub.clients) == 0 {
				logctx.Debug(ctx, "no active subscribers, dropping message", logger.String("channel", channel))
				r.mu.Unlock()
				continue
			}

			for clientChan := range sub.clients {
				select {
				case clientChan <- []byte(msg.Payload):
				default:
					logctx.Error(ctx, "client channel is full, dropping message", logger.String("channel", channel), logger.Int("buffer_size", getBufferSize()), logger.Int("clients", len(sub.clients)), logger.Int("buffered_messages", len(sub.messages)))
				}
			}
			r.mu.Unlock()
		case <-ctx.Done():
			logctx.Debug(ctx, "context cancelled", logger.String("channel", channel))
			return
		}
	}
}

// getBufferSize returns the configurable buffer size for the channelSubscription messages channel
func getBufferSize() int {
	if BUFFER_SIZE == "" {
		return 100
	}
	bufferSize, err := strconv.Atoi(BUFFER_SIZE)
	if err != nil {
		return 100
	}
	return bufferSize
}
