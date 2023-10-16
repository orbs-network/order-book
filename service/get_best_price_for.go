package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error) {
	order, err := s.orderBookStore.GetBestPriceFor(ctx, symbol, side)

	if err == models.ErrOrderNotFound {
		logctx.Info(ctx, fmt.Sprintf("No orders found for %s %q", side, symbol))
		return decimal.Zero, nil
	}

	if err != nil {
		logctx.Warn(ctx, fmt.Sprintf("Failed to get best price for %s %q: %s", side, symbol, err))
		return decimal.Zero, err
	}

	logctx.Info(ctx, fmt.Sprintf("Best price for %s %q is %q", side, symbol, order.Price))
	return order.Price, nil
}
