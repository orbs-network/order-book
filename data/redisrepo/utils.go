package redisrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/redis/go-redis/v9"
)

// CreateUserOpenOrdersKey creates a Redis key for storing the user's open orders
func CreateUserOpenOrdersKey(userId uuid.UUID, symbol models.Symbol) string {
	return fmt.Sprintf("userId:%s:%s:openOrders", userId, symbol.String())
}

// CreateOrderIDKey creates a Redis key for a single order ID
func CreateOrderIDKey(orderId uuid.UUID) string {
	return fmt.Sprintf("orderID:%s:order", orderId)
}
func CreateClosedOrderKey(orderId uuid.UUID) string {
	return fmt.Sprintf("order:closed:%s", orderId)
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

// CreateOpenSwapKey creates a Redis key for storing the swap data
func CreateOpenSwapKey(swapId uuid.UUID) string {
	return fmt.Sprintf("swap:open:%s", swapId)
}

func CreateResolvedSwapKey(swapId uuid.UUID) string {
	return fmt.Sprintf("swap:resolved:%s", swapId)
}

// CreateUserApiKeyKey creates a Redis key for storing the user by their API key
func CreateUserApiKeyKey(apiKey string) string {
	return fmt.Sprintf("userApiKey:%s:user", apiKey)
}

// CreateUserIdKey creates a Redis key for storing the user by their ID
func CreateUserIdKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:user", userId)
}

// CreateCompletedSwapsKey creates a Redis key for storing completed swaps
func CreateUserResolvedSwapsKey(userId uuid.UUID) string {
	return fmt.Sprintf("userId:%s:resolvedSwaps", userId)
}

// GENERIC store funcs
func AddVal2Set(ctx context.Context, client redis.Cmdable, key, val string) error {
	added, err := client.SAdd(ctx, key, val).Result()
	if err != nil {
		logctx.Error(ctx, "Failed to add element to set", logger.Error(err), logger.String("key", key), logger.String("val", val))
		return err
	}

	if added == 0 {
		logctx.Warn(ctx, "Element already in set", logger.String("key", key), logger.String("val", val))
		return models.ErrValAlreadyInSet
	}

	logctx.Debug(ctx, "Added element to set", logger.String("key", key), logger.String("val", val))
	return nil
}

func GetMakerTokenTrackKey(token, wallet string) string {
	return fmt.Sprintf("balance:%s:%s", strings.ToUpper(token), strings.ToUpper(wallet))
}

// key for tracking balance of a maker in a certain token
func Order2MakerTokenTrackKey(order models.Order) string {
	if order.Signature.AbiFragment.Info.Swapper.String() == "" {
		fmt.Println("order does not have a swapper address in ABI", logger.String("orderId", order.Id.String()))
		return ""
	}
	if order.Signature.AbiFragment.Input.Token.String() == "" {
		fmt.Println("order does not have an Input token address address in ABI", logger.String("orderId", order.Id.String()))
		return ""
	}
	return GetMakerTokenTrackKey(order.Signature.AbiFragment.Input.Token.String(), order.Signature.AbiFragment.Info.Swapper.String())
}
