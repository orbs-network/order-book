package redisrepo

import (
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func CreatePriceKey(symbol models.Symbol, price decimal.Decimal) string {
	return fmt.Sprintf("%s:orders:%s", symbol, price)
}

func CreateOrderIDKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:orderIDs", symbol)
}

func CreateBuySidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:buy:prices", symbol)
}

func CreateSellSidePricesKey(symbol models.Symbol) string {
	return fmt.Sprintf("%s:sell:prices", symbol)
}
