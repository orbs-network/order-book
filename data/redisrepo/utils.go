package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
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

// CreateUserPubKeyKey creates a Redis key for storing the user's public key
func CreateUserPubKeyKey(publicKey string) string {
	return fmt.Sprintf("user:%s:publicKey", publicKey)
}

// CreateAuctionTrackerKey creates a Redis key for storing auctions of different statuses
func CreateAuctionTrackerKey(status models.AuctionStatus) string {
	return fmt.Sprintf("auctionTracker:%s", status)
}

// GENERIC store funcs
func AddVal2Set(ctx context.Context, client redis.Cmdable, key, val string) error {
	added, err := client.SAdd(ctx, key, val).Result()
	if err != nil {
		return err
	}

	if added == 0 {
		return models.ErrValAlreadyInSet
	}

	return nil
}
