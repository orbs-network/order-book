package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (r *redisRepository) GetOrdersAtPrice(ctx context.Context, symbol models.Symbol, price decimal.Decimal) ([]models.Order, error) {
	panic("not implemented")
}
