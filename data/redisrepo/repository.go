package redisrepo

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client  redis.Cmdable
	txMap   map[uint]redis.Pipeliner
	ixIndex uint
}

func NewRedisRepository(client redis.Cmdable) (*redisRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client cannot be nil")
	}
	txMap := make(map[uint]redis.Pipeliner)
	return &redisRepository{
		client: client,
		txMap:  txMap,
	}, nil
}
