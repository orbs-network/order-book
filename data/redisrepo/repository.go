package redisrepo

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	cmdable       redis.Cmdable
	client        *redis.Client
	txMap         map[uint]redis.Pipeliner
	ixIndex       uint
	subscriptions map[string]*channelSubscription
	mu            sync.Mutex
}

type channelSubscription struct {
	pubsub   *redis.PubSub
	clients  map[chan []byte]struct{}
	messages chan []byte
}

func NewRedisRepository(cmdable redis.Cmdable) (*redisRepository, error) {
	if cmdable == nil {
		return nil, fmt.Errorf("redis client cannot be nil")
	}

	client, ok := cmdable.(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("cmdable is not a *redis.Client")
	}

	txMap := make(map[uint]redis.Pipeliner)
	return &redisRepository{
		cmdable:       cmdable,
		client:        client,
		txMap:         txMap,
		subscriptions: make(map[string]*channelSubscription),
	}, nil
}
