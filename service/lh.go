package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type Order interface {
	GetMinPriceOrder() *models.Order
	GetNextOrder() *models.Order
}

func (s *Service) GetAmountOut(ctx context.Context, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (decimal.Decimal, error) {

	var it OrderIter
	if side == models.SELL {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
	} else {
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
	}

	// buy 2 eth for 2000 usd
	// amount In = 2000, price = 1000
	amountOut := decimal.NewFromInt(0)
	var order *models.Order
	for it.HasNext() && amountIn.IsPositive() {
		order = it.Next()

		// max buy
		maxBuy := amountIn.Div(order.Price)

		minSize := decimal.Min(maxBuy, order.Size)

		amountIn.Sub(maxBuy)
		amountOut.Add(minSize)
	}
	return amountOut, nil
}
