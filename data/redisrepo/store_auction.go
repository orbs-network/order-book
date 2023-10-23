package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

func (r *redisRepository) StoreAuction(ctx context.Context, auctionID string, fillOrders []models.FilledOrder) error {
	// auctionId:<ID>: [{orderID: <ID>, filledAmount: <amount>}, ...}]
	panic("not implemented")
}
