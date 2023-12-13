package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
)

// CreateUserOpenOrdersKey creates a Redis key for storing the user's open orders
func CreateUserOpenOrdersKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:openOrders", userId)
}

// CreateUserFilledOrdersKey creates a Redis key for storing the user's filled orders
func CreateUserFilledOrdersKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:filledOrders", userId)
}

// CreateOrderIDKey creates a Redis key for a single order ID
func CreateOrderIDKey(orderId uuid.UUID) string {
	return fmt.Sprintf("orderID:%s:order", orderId)
}

// CreateClientOIDKey creates a Redis key for a single client order ID
func CreateClientOIDKey(clientOId uuid.UUID) string {
	return fmt.Sprintf("clientOId:%s:order", clientOId)
}

// CreateBuySidePricesKey creates a Redis key for storing the buy side (bid) prices
func CreateBuySidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:buy:prices", symbol)
}

// CreateSellSidePricesKey creates a Redis key for storing the sell side (ask) prices
func CreateSellSidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:sell:prices", symbol)
}

// CreateSwapKey creates a Redis key for storing the swap data
func CreateSwapKey(swapId uuid.UUID) string {
	return fmt.Sprintf("swapId:%s", swapId)
}

// CreateUserApiKeyKey creates a Redis key for storing the user by their API key
func CreateUserApiKeyKey(apiKey string) string {
	return fmt.Sprintf("userApiKey:%s:user", apiKey)
}

// CreateUserIdKey creates a Redis key for storing the user by their ID
func CreateUserIdKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:user", userId)
}

// CreateSwapTrackerKey creates a Redis key for storing swaps of different statuses
func CreateSwapTrackerKey(status models.SwapStatus) string {
	return fmt.Sprintf("swapTracker:%s", status)
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
