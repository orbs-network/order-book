package redisrepo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func CreatePriceKey(symbol models.Symbol, price decimal.Decimal) string {
	return fmt.Sprintf("%s:bestPrice:%s", symbol, price)
}

func CreateOrderIDKey(orderId uuid.UUID) string {
	return fmt.Sprintf("%s:orderID", orderId)
}

func CreateBuySidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:buy:prices", symbol)
}

func CreateSellSidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:sell:prices", symbol)
}
