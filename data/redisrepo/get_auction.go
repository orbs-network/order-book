package redisrepo

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) GetAuction(ctx context.Context, auctionID uuid.UUID) ([]models.FilledOrder, error) {
	auctionKey := CreateAuctionKey(auctionID)

	auctionJsons, err := r.client.LRange(ctx, auctionKey, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
		return []models.FilledOrder{}, models.ErrUnexpectedError
	}

	var allFilledOrders []models.FilledOrder

	for _, auctionJson := range auctionJsons {
		var orders []models.FilledOrder
		err := json.Unmarshal([]byte(auctionJson), &orders)
		if err != nil {
			logctx.Error(ctx, "failed to unmarshal auction", logger.String("auctionID", auctionID.String()), logger.Error(err))
			return []models.FilledOrder{}, models.ErrMarshalError
		}
		allFilledOrders = append(allFilledOrders, orders...)
	}

	logctx.Info(ctx, "got auction", logger.String("auctionID", auctionID.String()))
	return allFilledOrders, nil
}
