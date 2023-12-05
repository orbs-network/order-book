package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) BeginSwap(ctx context.Context, data models.AmountOut) (models.BeginSwapRes, error) {
	// create swapID
	swapId := uuid.New()
	// no re-entry is needed
	// err := s.orderBookStore.UpdateAuctionTracker(ctx, models.swap_started, auctionId)

	// storeSwap

	res := models.BeginSwapRes{
		OutAmount: data.Size,
		SwapId:    swapId,
	}

	// validate all orders of auction
	for _, frag := range data.OrderFrags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		if err != nil {
			logctx.Warn(ctx, err.Error())
			return models.BeginSwapRes{}, models.ErrOrderNotFound
		} else if !validateOrderFrag(frag, order) {
			// cancel auction
			_ = s.orderBookStore.RemoveSwap(ctx, swapId)

			// return empty
			logctx.Warn(ctx, "failed to validate order frag")
			return models.BeginSwapRes{}, models.ErrAuctionInvalid
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
	err := s.orderBookStore.StoreOrders(ctx, res.Orders)
	if err != nil {
		logctx.Error(ctx, "StoreOrders Failed", logger.Error(err))
		return models.BeginSwapRes{}, err
	}

	// add oredebook signature on the buffer
	//res.BookSignature = []byte("todo:sign")

	return res, nil
}

func (m *Service) AbortSwap(ctx context.Context, swapId uuid.UUID) error {
	return nil
}

// func (m *Service) txSent(ctx context.Context, swapId uuid.UUID) error {
// 	return nil
// }
