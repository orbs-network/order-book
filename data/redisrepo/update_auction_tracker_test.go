package redisrepo

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepository_UpdateAuctionTrackerTest(t *testing.T) {

	auctionStatus := models.AUCTION_REVERTED
	auctionId := uuid.MustParse("00000000-0000-0000-0000-000000000009")

	t.Run("should return without error when a new element is added that does not already exist", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateAuctionTrackerKey(auctionStatus)

		mock.ExpectSAdd(key, auctionId.String()).SetVal(1)

		err := repo.UpdateAuctionTracker(ctx, auctionStatus, auctionId)

		assert.NoError(t, err)
	})

	t.Run("should return `ErrValAlreadyInSet` error when element already exists", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateAuctionTrackerKey(auctionStatus)

		mock.ExpectSAdd(key, auctionId.String()).SetVal(0)

		err := repo.UpdateAuctionTracker(ctx, auctionStatus, auctionId)
		assert.ErrorIs(t, err, models.ErrValAlreadyInSet)
	})

	t.Run("should return error when unexpected Redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		repo := &redisRepository{
			client: db,
		}

		key := CreateAuctionTrackerKey(auctionStatus)

		mock.ExpectSAdd(key, auctionId.String()).SetErr(assert.AnError)

		err := repo.UpdateAuctionTracker(ctx, auctionStatus, auctionId)
		assert.ErrorContains(t, err, "failed to add auction to tracker")
	})
}
