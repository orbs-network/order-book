package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// update swap fields
//
// move swap to resolved key
//
// update orders state
//
// save swap to the users involved
func (e *EvmClient) ResolveSwap(ctx context.Context, swap models.Swap, isSuccessful bool, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	// resolve date
	swap.Resolved = time.Now()

	// success status
	swap.Succeeded = isSuccessful

	// save to "swapResolved" key
	// remove from active "swapId"
	err := e.orderBookStore.ResolveSwap(ctx, swap)
	if err != nil {
		logctx.Error(ctx, "Failed to ResolveSwap in store", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return err
	}

	// Failed     ===========================================================
	// same impl as abort swap
	if !isSuccessful {
		// unlock orders
		// mutual impl for ABORT and RESOLVE(false) swap
		err := unlockSwapAndHandleCancelledOrders(ctx, nil, e.orderBookStore, &swap)
		if err != nil {
			logctx.Error(ctx, "Failed unlockSwapAndHandleCancelledOrders", logger.Error(err), logger.String("swapId", swap.Id.String()))
		}
		return err
	}

	// successful ===========================================================
	var orderIds []uuid.UUID
	for _, frag := range swap.Frags {
		orderIds = append(orderIds, frag.OrderId)
	}

	// get orders from frags
	orders, err := e.orderBookStore.FindOrdersByIds(ctx, orderIds, false)
	if err != nil {
		logctx.Error(ctx, "Failed to get orders", logger.Error(err), logger.String("swapId", swap.Id.String()))
		return fmt.Errorf("failed to get orders: %w", err)
	}

	// get users from orders
	userIds := make(map[uuid.UUID]bool)
	filledOrders := []models.Order{}
	updatedOrders := []models.Order{}

	for i, order := range orders {
		// fill part of the order
		if _, err := order.Fill(ctx, swap.Frags[i]); err != nil {
			logctx.Error(ctx, "Failed to mark order as filled", logger.Error(err), logger.String("orderId", order.Id.String()))
			continue
		}

		// publish Fill Event
		fill := models.NewFill(order.Symbol, swap, swap.Frags[i], &order)
		e.publishFillEvent(ctx, order.UserId, *fill)

		if order.IsFilled() {
			// add to filled orders if completely filled
			filledOrders = append(filledOrders, order)
		} else {
			// update fill status
			updatedOrders = append(updatedOrders, order)
		}

		e.publishOrderEvent(ctx, &order)

		// append to users to be updated later
		userIds[order.UserId] = true
	}

	// update user(s) keys
	for userId := range userIds {
		// save resolved swap to a user
		err = e.orderBookStore.StoreUserResolvedSwap(ctx, userId, swap)
		if err != nil {
			logctx.Error(ctx, "Error StoreUserResolvedSwap", logger.Error(err), logger.String("swapId", swap.Id.String()))
		}
	}
	// save COMPLETELY filled orders in case of success and fill
	err = e.orderBookStore.StoreFilledOrders(ctx, filledOrders)
	if err != nil {
		logctx.Error(ctx, "Error StoreFilledOrders", logger.Error(err), logger.String("swapId", swap.Id.String()))
	}
	// save updated orders
	err = e.orderBookStore.StoreOpenOrders(ctx, updatedOrders)
	if err != nil {
		logctx.Error(ctx, "Error StoreOpenOrders", logger.Error(err), logger.String("swapId", swap.Id.String()))
	}

	logctx.Debug(ctx, "Resolved swap", logger.String("swapId", swap.Id.String()), logger.Bool("isSuccessful", isSuccessful), logger.String("created", swap.Created.String()), logger.String("resolved", swap.Resolved.String()), logger.String("txHash", swap.TxHash))
	return nil
}
