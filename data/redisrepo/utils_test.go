package redisrepo

import (
	"context"
	"fmt"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_Utils(t *testing.T) {

	// no mock!!!
	address := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "secret",
		DB:       10, // for test
	})
	if rdb == nil {
		fmt.Println("redis is not running in ", address)
		return
	}

	repository, err := NewRedisRepository(rdb)
	assert.NoError(t, err, "should not return error")

	t.Run("AddVal2Set  - should set userfail if value already exist", func(t *testing.T) {
		const val1 = "val1"
		const val2 = "val2"
		const key = "test-set"
		ctx := context.Background()
		err := repository.AddVal2Set(ctx, key, val1)
		assert.NoError(t, err, "should not return error")
		err = repository.AddVal2Set(ctx, key, val1)
		assert.EqualError(t, err, models.ErrValAlreadyInSet.Error())
		err = repository.AddVal2Set(ctx, key, val2)
		assert.NoError(t, err, "should not return error")
		repository.client.Del(ctx, key)
	})
}
