package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// CreateOrderIDKey creates a Redis key for storing the user's orders
func CreateUserOrdersKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:orders", userId)
}

// CreateOrderIDKey creates a Redis key for a single order ID
func CreateOrderIDKey(orderId uuid.UUID) string {
	return fmt.Sprintf("orderID:%s", orderId)
}

// CreateClientOIDKey creates a Redis key for a single client order ID
func CreateClientOIDKey(clientOId uuid.UUID) string {
	return fmt.Sprintf("clientOId:%s", clientOId)
}

// CreateBuySidePricesKey creates a Redis key for storing the buy side (bid) prices
func CreateBuySidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:buy:prices", symbol)
}

// CreateSellSidePricesKey creates a Redis key for storing the sell side (ask) prices
func CreateSellSidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:sell:prices", symbol)
}

// CreateAuctionKey creates a Redis key for storing the auction data
func CreateAuctionKey(auctionID uuid.UUID) string {
	return fmt.Sprintf("auctionId:%s", auctionID)
}

// GENERIC store funcs
func (r *redisRepository) AddVal2Set(ctx context.Context, key, val string) error {
	isMember, err := r.client.SIsMember(ctx, key, val).Result()
	if err != nil {
		logctx.Warn(ctx, "SIsMember Failed", logger.Error(err))
		return err
	}
	if isMember {
		return models.ErrValAlreadyInSet
	}

	_, err = r.client.SAdd(ctx, key, val).Result()
	if err != nil {
		logctx.Warn(ctx, "SAdd Failed", logger.Error(err))
		return err
	}

	return nil
}
