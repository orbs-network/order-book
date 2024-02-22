package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (e *EvmClient) ProcessCompletedTransaction(ctx context.Context, tx *models.Tx, swapId uuid.UUID, isSuccessful bool, mu *sync.Mutex) ([]models.Order, error) {
	mu.Lock()
	defer mu.Unlock()

	swap, err := e.orderBookStore.GetSwap(ctx, swapId)
	if err != nil {
		logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapId", swapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get swap: %w", err)
	}

	var orderIds []uuid.UUID
	orderSizes := make(map[uuid.UUID]decimal.Decimal)

	for _, frag := range swap.Frags {
		orderIds = append(orderIds, frag.OrderId)
		orderSizes[frag.OrderId] = frag.OutSize
	}

	orders, err := e.orderBookStore.FindOrdersByIds(ctx, orderIds, false)
	if err != nil {
		logctx.Error(ctx, "Failed to get orders", logger.Error(err), logger.String("swapId", swapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get orders: %w", err)
	}

	var swapOrdersWithSize []store.OrderWithSize

	for _, order := range orders {

		size, found := orderSizes[order.Id]
		if !found {
			logctx.Error(ctx, "Failed to get order frag size", logger.String("orderId", order.Id.String()))
			return []models.Order{}, fmt.Errorf("failed to get order frag size")
		}

		if isSuccessful {
			if _, err := order.Fill(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to mark order as filled", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}
		} else {
			if err := order.Unlock(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to Release order locked liq", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}
		}
		swapOrdersWithSize = append(swapOrdersWithSize, store.OrderWithSize{
			Order: &order,
			Size:  size,
		})
	}

	err = e.orderBookStore.ProcessCompletedSwapOrders(ctx, swapOrdersWithSize, swapId, tx, isSuccessful)
	if err != nil {
		logctx.Error(ctx, "Failed to process completed swap orders", logger.Error(err), logger.String("swapId", swapId.String()))
		return []models.Order{}, fmt.Errorf("failed to process completed swap orders: %w", err)
	}

	return orders, nil
}
