package redisrepo

import (
	"context"
	"fmt"
	"sync"
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

	// Fetch sell and buy side orders concurrently
	var askOrderIds, bidOrderIds []uuid.UUID
	var asksErr, bidsErr error
	var wg sync.WaitGroup
	wg.Add(2)

	// Fetch Asks
	go func() {
		defer wg.Done()
		askOrderIds, asksErr = r.fetchOrderIds(ctx, CreateSellSidePricesKey(symbol), depth)
	}()

	// Fetch Bids
	go func() {
		defer wg.Done()
		bidOrderIds, bidsErr = r.fetchOrderIds(ctx, CreateBuySidePricesKey(symbol), depth)
	}()

	wg.Wait()
	if asksErr != nil {
		logctx.Error(ctx, "failed to fetch asks", logger.Error(asksErr))
		return marketDepth, fmt.Errorf("failed to fetch asks: %v", asksErr)
	}
	if bidsErr != nil {
		logctx.Error(ctx, "failed to fetch bids", logger.Error(bidsErr))
		return marketDepth, fmt.Errorf("failed to fetch bids: %v", bidsErr)
	}

	// Concatenate ask and bid order IDs
	allOrderIds := append(askOrderIds, bidOrderIds...)

	// Fetch all orders in bulk
	var orders []models.Order
	var err error
	for i := 0; i < len(allOrderIds); i += MAX_ORDER_IDS {
		end := i + MAX_ORDER_IDS
		if end > len(allOrderIds) {
			end = len(allOrderIds)
		}

		orders, err = r.FindOrdersByIds(ctx, allOrderIds[i:end], false)
		if err != nil {
			logctx.Error(ctx, "failed to find orders by IDs", logger.Error(err))
			return marketDepth, fmt.Errorf("failed to find orders by IDs: %v", err)
		}
	}

	// Map orders by their ID for quick lookup
	orderMap := make(map[uuid.UUID]models.Order)
	for _, order := range orders {
		orderMap[order.Id] = order
	}

	// Process asks
	for _, orderId := range askOrderIds {
		if order, ok := orderMap[orderId]; ok {
			marketDepth.Asks = append(marketDepth.Asks, []decimal.Decimal{order.Price, order.GetAvailableSize()})
		}
	}

	// Process bids
	for _, orderId := range bidOrderIds {
		if order, ok := orderMap[orderId]; ok {
			marketDepth.Bids = append(marketDepth.Bids, []decimal.Decimal{order.Price, order.GetAvailableSize()})
		}
	}

	logctx.Info(ctx, "fetched market depth", logger.String("symbol", symbol.String()), logger.Int("depth", depth), logger.Int("numAsks", len(marketDepth.Asks)), logger.Int("numBids", len(marketDepth.Bids)))
	return marketDepth, nil
}

// fetchOrderIds fetches order IDs for a given key and depth
func (r *redisRepository) fetchOrderIds(ctx context.Context, key string, depth int) ([]uuid.UUID, error) {
	orderStrIds, err := r.client.ZRange(ctx, key, 0, int64(depth-1)).Result()
	if err != nil {
		logctx.Error(ctx, "failed to fetch order IDs from Redis", logger.Error(err), logger.String("key", key), logger.Int("depth", depth))
		return nil, fmt.Errorf("failed to fetch order IDs from Redis: %v", err)
	}

	orderIds := make([]uuid.UUID, len(orderStrIds))
	for i, strId := range orderStrIds {
		orderId, err := uuid.Parse(strId)
		if err != nil {
			logctx.Error(ctx, "failed to parse order ID", logger.Error(err), logger.String("key", key), logger.Int("depth", depth), logger.String("strId", strId))
			return nil, fmt.Errorf("failed to parse order ID: %v", err)
		}
		orderIds[i] = orderId
	}

	logctx.Info(ctx, "fetched order IDs", logger.String("key", key), logger.Int("depth", depth), logger.Int("numOrders", len(orderIds)))
	return orderIds, nil
}
