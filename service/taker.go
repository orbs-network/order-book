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
	return order.GetAvailableSize().GreaterThanOrEqual(frag.OutSize)
}

func validatePendingFrag(frag models.OrderFrag, order *models.Order) bool {
	// check if order is still open
	if order.IsFilled() {
		return false
	}
	// order.Size pending should be greater or equal to orderFrag: (Order.sizePending + Order.pending) >= frag.size
	return order.SizePending.GreaterThanOrEqual(frag.OutSize)
}

func (s *Service) BeginSwap(ctx context.Context, data models.QuoteRes) (models.BeginSwapRes, error) {
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
		err := res.Orders[i].Lock(ctx, res.Fragments[i].OutSize)
		if err != nil {
			logctx.Error(ctx, "Lock order Failed", logger.Error(err))
			return models.BeginSwapRes{}, err
		}
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

	// add oredebook signature on the buffer HERE if needed
	return res, nil
}

func (s *Service) SwapStarted(ctx context.Context, swapId uuid.UUID, txHash string) error {
	logctx.Debug(ctx, "SwapStarted", logger.String("swapId", swapId.String()))
	err := s.orderBookStore.StoreNewPendingSwap(ctx, models.SwapTx{
		SwapId: swapId,
		TxHash: txHash,
	})
	if err != nil {
		logctx.Error(ctx, "StoreNewPendingSwap failed", logger.Error(err))
	}
	return err
}

func (s *Service) AbortSwap(ctx context.Context, swapId uuid.UUID) error {
	logctx.Info(ctx, "AbortSwap", logger.String("swapId", swapId.String()))
	// get swap from store
	swap, err := s.orderBookStore.GetSwap(ctx, swapId)
	if err != nil {
		logctx.Warn(ctx, "GetSwap Failed", logger.Error(err))
		return err
	}

	orders := []models.Order{}
	// validate all pending orders fragments of auction
	for _, frag := range swap.Frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		// no return during erros as what can be revert, should
		if err != nil {
			logctx.Error(ctx, "order not found while reverting a swap", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Swap fragments should be valid during a revert request", logger.Error(err))
		} else {
			// success
			logctx.Debug(ctx, "Unlock Fragment", logger.String("orderID", frag.OrderId.String()), logger.String("OutSize", frag.OutSize.String()))
			err = order.Unlock(ctx, frag.OutSize)
			if err != nil {
				logctx.Error(ctx, "Unlock Failed", logger.Error(err))
				return err
			}
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

func (s *Service) FillSwap(ctx context.Context, swapId uuid.UUID) error {
	logctx.Debug(ctx, "FillSwap", logger.String("swapId", swapId.String()))

	// get swap from store
	swap, err := s.orderBookStore.GetSwap(ctx, swapId)
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
			logctx.Error(ctx, "order not found while reverting a swap", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Swap fragments should be valid during a revert request", logger.Error(err))
		} else {
			// from pending to fill
			logctx.Debug(ctx, "FillOrder", logger.String("orderID", order.Id.String()), logger.String("OutSize", frag.OutSize.String()))
			filled, err := order.Fill(ctx, frag.OutSize)
			if err != nil {
				logctx.Error(ctx, "FillOrder Failed", logger.Error(err))
				return err
			}
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
