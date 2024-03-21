package redisrepo

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	cmdable redis.Cmdable
	client  *redis.Client
	txMap   map[uint]redis.Pipeliner
	ixIndex uint
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
		cmdable: cmdable,
		client:  client,
		txMap:   txMap,
	}, nil
}
