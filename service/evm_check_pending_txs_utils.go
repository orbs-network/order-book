package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// update swap fields
// move swap to resolved key
// update orders state
// save swap to the users involved

func (e *EvmClient) ResolveSwap(ctx context.Context, swap models.Swap, isSuccessful bool, mu *sync.Mutex) ([]models.Order, error) {
	mu.Lock()
	defer mu.Unlock()

	//swap, err := e.orderBookStore.GetSwap(ctx, swap.Id)
	// if err != nil {
	// 	logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapId", swap.Id))
	// 	return []models.Order{}, fmt.Errorf("failed to get swap: %w", err)
	// }

	// resolve date
	swap.Resolved = time.Now()

	// success status
	swap.Succeeded = isSuccessful

	// save to "swapResolved" key
	// remove from active "swapId"
	err := e.orderBookStore.ResolveSwap(ctx, swap)

	if err != nil {
		logctx.Error(ctx, "Failed to ResolveSwap in store", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return []models.Order{}, err
	}

	var orderIds []uuid.UUID
	orderSizes := make(map[uuid.UUID]decimal.Decimal)

	for _, frag := range swap.Frags {
		orderIds = append(orderIds, frag.OrderId)
		orderSizes[frag.OrderId] = frag.OutSize
	}

	// get orders from frags
	orders, err := e.orderBookStore.FindOrdersByIds(ctx, orderIds, false)
	if err != nil {
		logctx.Error(ctx, "Failed to get orders", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return []models.Order{}, fmt.Errorf("failed to get orders: %w", err)
	}

	//var swapOrdersWithSize []store.OrderWithSize

	// get users from orders
	userIds := make(map[uuid.UUID]bool)
	filledOrders := []models.Order{}
	updatedOrders := []models.Order{}

	for _, order := range orders {
		size, found := orderSizes[order.Id]
		if !found {
			logctx.Error(ctx, "Failed to get order frag size", logger.String("orderId", order.Id.String()))
			return []models.Order{}, fmt.Errorf("failed to get order frag size")
		}

		if isSuccessful {
			// fill part of the order
			if _, err := order.Fill(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to mark order as filled", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}

			if order.IsFilled() {
				// add to filled orders if completely filled
				filledOrders = append(filledOrders, order)
			} else {
				// update fill status
				updatedOrders = append(updatedOrders, order)
			}
		} else {
			// unlock orders
			if err := order.Unlock(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to Release order locked liq", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}
			updatedOrders = append(updatedOrders, order)
		}

		// append to users
		userIds[order.UserId] = true
		// swapOrdersWithSize = append(swapOrdersWithSize, store.OrderWithSize{
		// 	Order: &order,
		// 	Size:  size,
		// })
	}
	// TODO: save orders new state

	// update user(s) keys
	for userId := range userIds {
		// save resolved swap to a user
		err = e.orderBookStore.StoreUserResolvedSwap(ctx, userId, swap)
		if err != nil {
			logctx.Error(ctx, "Error StoreUserResolvedSwap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		}
	}
	// save fill orders in case of success and fill
	err = e.orderBookStore.StoreFilledOrders(ctx, filledOrders)
	if err != nil {
		logctx.Error(ctx, "Error StoreFilledOrders", logger.Error(err), logger.String("swapId", swap.Id.String()))
	}
	// save updated orders
	err = e.orderBookStore.StoreOpenOrders(ctx, updatedOrders)
	if err != nil {
		logctx.Error(ctx, "Error StoreOpenOrders", logger.Error(err), logger.String("swapId", swap.Id.String()))
	}

	// err = e.orderBookStore.ProcessCompletedSwapOrders(ctx, swapOrdersWithSize, swap.Id, swap.TxHash, isSuccessful)
	// if err != nil {
	// 	logctx.Error(ctx, "Failed to process completed swap orders", logger.Error(err), logger.String("swapId", swap.Id.String()))
	// 	return []models.Order{}, fmt.Errorf("failed to process completed swap orders: %w", err)
	// }

	return orders, nil

}

// replaced by ResolveSwap
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
