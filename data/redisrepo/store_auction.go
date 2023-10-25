package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) StoreAuction(ctx context.Context, auctionID uuid.UUID, fillOrders []models.FilledOrder) error {

	auctionJson, err := models.MarshalFilledOrders(fillOrders)
	if err != nil {
		logctx.Error(ctx, "failed to marshal auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
		return models.ErrMarshalError
	}

	auctionKey := CreateAuctionKey(auctionID)

	_, err = r.client.RPush(ctx, auctionKey, auctionJson).Result()
	if err != nil {
		logctx.Error(ctx, "failed to store auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
		return models.ErrUnexpectedError
	}

	logctx.Info(ctx, "stored auction", logger.String("auctionID", auctionID.String()))
	return nil
}
