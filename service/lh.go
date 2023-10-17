package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (s *Service) GetAmountOut(ctx context.Context, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (decimal.Decimal, error) {
	// get
	//order, err := s.orderBookStore.GetBestPriceFor(ctx, symbol, side)
	//order.

	// if err == models.ErrOrderNotFound {
	// 	logctx.Info(ctx, fmt.Sprintf("No orders found for %s %q", side, symbol))
	// 	return decimal.Zero, nil
	// }

	// if err != nil {
	// 	logctx.Warn(ctx, fmt.Sprintf("Failed to get best price for %s %q: %s", side, symbol, err))
	// 	return decimal.Zero, err
	// }

	// logctx.Info(ctx, fmt.Sprintf("Best price for %s %q is %q", side, symbol, order.Price))
	// return order.Price, nil

	return decimal.NewFromString("1111")
	// res, err := decimal.NewFromString("1111")
	// return res, nil
}
