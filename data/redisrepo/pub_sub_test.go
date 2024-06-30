package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var orderJson, _ = test_order.ToJson()

func TestRedisRepository_PublishEvent(t *testing.T) {

	t.Run("should publish event", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectPublish(test_order.Id.String(), orderJson).SetVal(1)

		err := repo.PublishEvent(ctx, test_order.Id.String(), orderJson)

		assert.NoError(t, err)
	})

	t.Run("should return error when failed to publish event", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		mock.ExpectPublish(test_order.Id.String(), orderJson).SetErr(assert.AnError)

		err := repo.PublishEvent(ctx, test_order.Id.String(), orderJson)

		assert.ErrorContains(t, err, "failed to publish redis event")
	})
}
