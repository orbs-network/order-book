package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func validateOrderFrag(frag models.OrderFrag, order *models.Order) bool {

	// check if order is still open
	// TODO: we no longer need this check as we are not storing open orders?
	if order.IsFilled() {
		return false
	}
	// order.size - (Order.filled + Order.pending) >= frag.size
	// always in A token
	size := order.FragAtokenSize(frag)
	return order.GetAvailableSize().GreaterThanOrEqual(size)
}

func validatePendingFrag(frag models.OrderFrag, order *models.Order) bool {
	// check if order is still open
	if order.IsFilled() {
		return false
	}
	// order.Size pending should be greater or equal to orderFrag: (Order.sizePending + Order.pending) >= frag.size
	size := order.FragAtokenSize(frag)
	return order.SizePending.GreaterThanOrEqual(size)
}

func (s *Service) BeginSwap(ctx context.Context, data models.QuoteRes) (models.BeginSwapRes, error) {
	logctx.Info(ctx, "BeginSwap start", logger.String("data.size", data.Size.String()))

	// create swapID
	swapId := uuid.New()
	// no re-entry is needed

	res := models.BeginSwapRes{
		OutAmount: data.Size,
		SwapId:    swapId,
	}

	// validate all orders of a swap
	for _, frag := range data.OrderFrags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		if err != nil {
			logctx.Warn(ctx, err.Error())
			return models.BeginSwapRes{}, models.ErrNotFound
		} else if !validateOrderFrag(frag, order) {
			// cancel swap
			_ = s.orderBookStore.RemoveSwap(ctx, swapId)

			// return empty
			logctx.Warn(ctx, "failed to validate order frag")
			return models.BeginSwapRes{}, models.ErrSwapInvalid
		} else {
			// success- append
			res.Orders = append(res.Orders, *order)
			res.Fragments = append(res.Fragments, frag)
		}
	}
	// lock liquidity - Only after fragments were validated
	// set order fragments as Pending
	for i := 0; i < len(res.Orders); i++ {
		// lock frag.Amount as pending per order - no STATUS_PENDING is needed
		logctx.Debug(ctx, "Lock Fragment", logger.String("orderID", res.Orders[i].Id.String()), logger.String("OutSize", res.Fragments[i].OutSize.String()))

		err := res.Orders[i].Lock(ctx, res.Fragments[i])
		if err != nil {
			logctx.Error(ctx, "Lock order Failed", logger.Error(err))
			return models.BeginSwapRes{}, err
		}
		s.publishOrderEvent(ctx, &res.Orders[i])
	}

	// save
	err := s.orderBookStore.StoreOpenOrders(ctx, res.Orders)
	if err != nil {
		logctx.Error(ctx, "StoreOrders Failed", logger.Error(err))
		return models.BeginSwapRes{}, err
	}

	err = s.orderBookStore.StoreSwap(ctx, swapId, res.Fragments)
	if err != nil {
		logctx.Error(ctx, "StoreSwap Failed", logger.Error(err))
		return models.BeginSwapRes{}, err
	}

	logctx.Info(ctx, "BeginSwap end ok", logger.String("swapId", string(res.SwapId.String())))

	// add oredebook signature on the buffer HERE if needed
	return res, nil
}

func (s *Service) SwapStarted(ctx context.Context, swapId uuid.UUID, txHash string) error {
	logctx.Info(ctx, "SwapStarted", logger.String("swapId", swapId.String()))
	err := s.orderBookStore.StoreNewPendingSwap(ctx, models.SwapTx{
		SwapId: swapId,
		TxHash: txHash,
	})
	if err != nil {
		logctx.Error(ctx, "StoreNewPendingSwap failed", logger.Error(err))
	}
	return err
}

// to be reused by resolveSwap
func unlockSwapAndHandleCancelledOrders(ctx context.Context, svc *Service, store store.OrderBookStore, swap *models.Swap) error {
	unlockedOrders := []models.Order{}
	ordersToRemove := []models.Order{}
	// validate all pending orders fragments of auction
	for _, frag := range swap.Frags {
		// get order by ID
		order, err := store.FindOrderById(ctx, frag.OrderId, false)
		// no return during erros as what can be revert, should
		if err != nil {
			logctx.Error(ctx, "order not found while reverting a swap", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Swap fragments should be valid during a revert request", logger.Error(err))
		} else {
			// success - Ulock an order even if cancelled so the condition IsPending below would be False and the order will be removed
			logctx.Debug(ctx, "Unlock Fragment", logger.String("orderID", frag.OrderId.String()), logger.String("OutSize", frag.OutSize.String()))
			err = order.Unlock(ctx, frag)
			if err != nil {
				logctx.Error(ctx, "Unlock Failed", logger.Error(err))
				return err
			}
			// no need to publish nor update cancelled order
			if !order.Cancelled {
				if svc != nil {
					svc.publishOrderEvent(ctx, order)
				}
				// save to modify/update new pending state in db
				unlockedOrders = append(unlockedOrders, *order)
			}
		}
		// remove cancelled unfilled non pending unlocked orders
		if order != nil && order.Cancelled && !order.IsPartialFilled() && !order.IsPending() {
			ordersToRemove = append(ordersToRemove, *order)
		}
	}
	err := store.PerformTx(ctx, func(txid uint) error {
		// update unlocked orders
		for _, order := range unlockedOrders {
			if err := store.TxModifyOrder(ctx, txid, models.Update, order); err != nil {
				logctx.Error(ctx, "AbortSwap Failed updating unlocked order", logger.Error(err), logger.String("orderId", order.Id.String()))
				return err
			}
		}
		// remove order from all entries
		for _, order := range ordersToRemove {
			if err := store.TxRemoveOrder(ctx, txid, order); err != nil {
				logctx.Error(ctx, "TxRemoveOrder Failed", logger.Error(err), logger.String("id", order.Id.String()))
				return err
			}
		}
		return nil
	})
	return err

}

func (s *Service) AbortSwap(ctx context.Context, swapId uuid.UUID) error {
	logctx.Debug(ctx, "AbortSwap", logger.String("swapId", swapId.String()))
	// get swap from store
	swap, err := s.orderBookStore.GetSwap(ctx, swapId, true)

	if err != nil {
		logctx.Warn(ctx, "GetSwap Failed", logger.Error(err))
		return err
	}

	// mutual impl for ABORT and RESOLVE(false) swap
	err = unlockSwapAndHandleCancelledOrders(ctx, s, s.orderBookStore, swap)
	if err != nil {
		logctx.Warn(ctx, "unlockSwapAndHandleCancelledOrders Failed", logger.Error(err))
		return err
	}

	return s.orderBookStore.RemoveSwap(ctx, swapId)
}

func (s *Service) FillSwap(ctx context.Context, swapId uuid.UUID) error {
	logctx.Debug(ctx, "FillSwap", logger.String("swapId", swapId.String()))

	// get swap from store
	swap, err := s.orderBookStore.GetSwap(ctx, swapId, true)
	if err != nil {
		logctx.Warn(ctx, "GetSwap Failed", logger.Error(err))
		return err
	}

	filledOrders := []models.Order{}
	openOrders := []models.Order{}
	// validate all pending orders fragments of auction
	for _, frag := range swap.Frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		// no return during erros as what can be revert, should
		if err != nil {
			logctx.Error(ctx, "order not found Fill Swap", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Swap pending fragments should be valid", logger.Error(err))
		} else {
			// from pending to fill
			logctx.Debug(ctx, "FillOrder", logger.String("orderID", order.Id.String()), logger.String("OutSize", frag.OutSize.String()))
			filled, err := order.Fill(ctx, frag)
			if err != nil {
				logctx.Error(ctx, "FillOrder Failed", logger.Error(err))
				return err
			}
			// publish fill event
			s.publishFillEvent(ctx, order.UserId, *models.NewFill(order.Symbol, *swap, frag, order))

			if filled {
				filledOrders = append(filledOrders, *order)

			} else {
				openOrders = append(openOrders, *order)
			}
		}
	}
	// store partial orders
	err = s.orderBookStore.StoreOpenOrders(ctx, openOrders)
	if err != nil {
		logctx.Warn(ctx, "StoreOrders Failed", logger.Error(err))
		return err
	}
	// store filled orders
	err = s.orderBookStore.StoreFilledOrders(ctx, filledOrders)
	if err != nil {
		logctx.Warn(ctx, "StoreFilledOrders Failed", logger.Error(err))
		return err
	}

	return s.orderBookStore.RemoveSwap(ctx, swapId)
}
