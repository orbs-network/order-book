package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// orderID->amount bought or sold in A token always

type ConfirmAuctionRes struct {
	Orders        []*models.Order
	Fragments     []*models.OrderFrag
	BookSignature []byte
}

func validateOrderFrag(frag models.OrderFrag, order *models.Order) bool {

	// check if order is still open
	if order.Status != models.STATUS_OPEN {
		return false
	}
	// order.size - (Order.filled + prder.pending) >= frag.size
	orderLockedSum := order.SizeFilled.Sub(order.SizePending)
	return order.Size.Sub(orderLockedSum).GreaterThanOrEqual(frag.Size)
}

func (s *Service) ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (ConfirmAuctionRes, error) {
	// get auction from store
	frags, err := s.orderBookStore.GetAuction(ctx, auctionId)
	if err != nil {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return ConfirmAuctionRes{}, models.ErrInsufficientLiquity
	}

	res := ConfirmAuctionRes{}

	// validate all orders of auction
	for _, frag := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		if order == nil {
			// cancel auction
			s.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrOrderNotFound
		} else if !validateOrderFrag(frag, order) {
			// cancel auction
			s.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrAuctionInvalid
		} else {
			// success- append
			res.Orders = append(res.Orders, order)
			res.Fragments = append(res.Fragments, &frag)

			// later s.orderBookStore.FillOrder()

		}
	}
	// process all fill requests
	for i := 0; i < len(res.Orders); i++ {
		// lock frag.Amount as pending per order - no STATUS_PENDING is needed
		res.Orders[i].SizePending = res.Fragments[i].Size
	}
	s.orderBookStore.StoreOrders(ctx, res.Orders)

	// add oredebook signature on the buffer
	res.BookSignature = []byte("todo:sign")

	// set entire auction as pending ??
	//s.orderBookStore.RemoveAuction(auctionId)

	// error

	return res, nil
}

func (s *Service) RemoveAuction(ctx context.Context, auctionId uuid.UUID) error {
	return nil
}
