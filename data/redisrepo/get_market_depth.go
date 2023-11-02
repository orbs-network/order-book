package redisrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (r *redisRepository) GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error) {
	// Create a default MarketDepth
	marketDepth := models.MarketDepth{
		Asks:   [][]decimal.Decimal{},
		Bids:   [][]decimal.Decimal{},
		Symbol: symbol.String(),
		Time:   time.Now().Unix(),
	}

	errChan := make(chan error, 2)

	// Asks
	sellSideKey := CreateSellSidePricesKey(symbol)
	go func() {
		asks, err := r.client.ZRange(ctx, sellSideKey, 0, int64(depth-1)).Result()
		if err != nil {
			logctx.Error(ctx, "Error fetching asks", logger.Error(err))
			errChan <- err
			return
		}

		for _, orderIDStr := range asks {
			orderId, err := uuid.Parse(orderIDStr)
			if err != nil {
				logctx.Error(ctx, "Error parsing ask order id", logger.Error(err))
				continue
			}
			order, err := r.FindOrderById(ctx, orderId, false)
			if err != nil {
				logctx.Error(ctx, "Error fetching order", logger.Error(err))
				continue
			}
			marketDepth.Asks = append(marketDepth.Asks, []decimal.Decimal{order.Price, order.GetAvailableSize()})
		}
		errChan <- nil
	}()

	// Bids
	buySideKey := CreateBuySidePricesKey(symbol)
	go func() {
		bids, err := r.client.ZRevRange(ctx, buySideKey, 0, int64(depth-1)).Result()
		if err != nil {
			logctx.Error(ctx, "Error fetching bids", logger.Error(err))
			errChan <- err
			return
		}

		for _, orderIDStr := range bids {
			orderId, err := uuid.Parse(orderIDStr)
			if err != nil {
				logctx.Error(ctx, "Error parsing bid order id", logger.Error(err))
				continue
			}
			order, err := r.FindOrderById(ctx, orderId, false)
			if err != nil {
				logctx.Error(ctx, "Error fetching order", logger.Error(err))
				continue
			}
			marketDepth.Bids = append(marketDepth.Bids, []decimal.Decimal{order.Price, order.GetAvailableSize()})
		}
		errChan <- nil
	}()

	// Wait for both goroutines to finish
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return marketDepth, err
		}
	}

	return marketDepth, nil
}
