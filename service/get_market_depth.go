package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	marketDepth, err := s.orderBookStore.GetMarketDepth(ctx, symbol, depth)

	if err != nil {
		logctx.Error(ctx, "unexpected error getting market depth", logger.Error(err), logger.String("symbol", string(symbol)))
		return models.MarketDepth{}, err
	}

	logctx.Debug(ctx, "got market depth", logger.String("symbol", symbol.String()))
	return marketDepth, nil
}
