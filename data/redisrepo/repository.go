package redisrepo

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client redis.Cmdable
}

func NewRedisRepository(client redis.Cmdable) (*redisRepository, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client cannot be nil")
	}

	return &redisRepository{
		client: client,
	}, nil
}
