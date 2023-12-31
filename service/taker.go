package service

import (
	"context"

	"github.com/google/uuid"
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
	return order.GetAvailableSize().GreaterThanOrEqual(frag.Size)
}

func validatePendingFrag(frag models.OrderFrag, order *models.Order) bool {
	// check if order is still open
	if order.IsFilled() {
		return false
	}
	// order.Size pending should be greater or equal to orderFrag: (Order.sizePending + Order.pending) >= frag.size
	return order.SizePending.GreaterThanOrEqual(frag.Size)
}

func (s *Service) BeginSwap(ctx context.Context, data models.AmountOut) (models.BeginSwapRes, error) {
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
		res.Orders[i].SizePending = res.Fragments[i].Size
	}
	err := s.orderBookStore.StoreOpenOrders(ctx, res.Orders)
	if err != nil {
		logctx.Error(ctx, "StoreOrders Failed", logger.Error(err))
		return models.BeginSwapRes{}, err
	}

	// add oredebook signature on the buffer HERE if needed

	return res, nil
}

func (s *Service) AbortSwap(ctx context.Context, swapId uuid.UUID) error {
	// returns error if already confirmed
	err := s.orderBookStore.UpdateSwapTracker(ctx, models.SWAP_ABORDTED, swapId)

	if err != nil {
		if err == models.ErrValAlreadyInSet {
			logctx.Warn(ctx, "AbortSwap re-entry!", logger.String("swapId: ", swapId.String()))
		} else {
			logctx.Warn(ctx, "AbortSwap UpdateSwapTracker Failed", logger.String("swapId: ", swapId.String()), logger.Error(err))
		}
		return err
	}

	// get swap from store
	frags, err := s.orderBookStore.GetSwap(ctx, swapId)
	if err != nil {
		logctx.Warn(ctx, "GetSwap Failed", logger.Error(err))
		return err
	}

	orders := []models.Order{}
	// validate all pending orders fragments of auction
	for _, frag := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		// no return during erros as what can be revert, should
		if err != nil {
			logctx.Error(ctx, "order not found while reverting a swap", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Swap fragments should be valid during a revert request", logger.Error(err))
		} else {
			// success
			order.SizePending = order.SizePending.Sub(frag.Size)
			orders = append(orders, *order)
		}
	}
	// store orders
	err = s.orderBookStore.StoreOpenOrders(ctx, orders)
	if err != nil {
		logctx.Warn(ctx, "StoreOrders Failed", logger.Error(err))
		return err
	}

	return s.orderBookStore.RemoveSwap(ctx, swapId)
}
