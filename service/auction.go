package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
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

func validatePendingFrag(frag models.OrderFrag, order *models.Order) bool {
	// check if order is still open
	if order.Status != models.STATUS_OPEN {
		return false
	}
	// order.Size pending should be greater or equal to orderFrag: (Order.sizePending + prder.pending) >= frag.size
	return order.SizePending.GreaterThanOrEqual(frag.Size)
}

func (s *Service) ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (ConfirmAuctionRes, error) {
	// TODO: re-entrance validate it doesnt already confirmed

	// get auction from store
	frags, err := s.orderBookStore.GetAuction(ctx, auctionId)
	if err != nil {
		logctx.Warn(ctx, "GetAuction Failed", logger.Error(err))
		return ConfirmAuctionRes{}, err
	}

	res := ConfirmAuctionRes{}

	// validate all orders of auction
	for _, frag := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		if order == nil {
			// cancel auction
			s.orderBookStore.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrOrderNotFound
		} else if !validateOrderFrag(frag, order) {
			// cancel auction
			s.orderBookStore.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrAuctionInvalid
		} else {
			// success- append
			res.Orders = append(res.Orders, order)
			res.Fragments = append(res.Fragments, &frag)
		}
	}
	// set order order fragments as Pending
	for i := 0; i < len(res.Orders); i++ {
		// lock frag.Amount as pending per order - no STATUS_PENDING is needed
		res.Orders[i].SizePending = res.Fragments[i].Size
	}
	s.orderBookStore.StoreOrders(ctx, res.Orders)

	// add oredebook signature on the buffer
	res.BookSignature = []byte("todo:sign")

	return res, nil
}

func (s *Service) RevertAuction(ctx context.Context, auctionId uuid.UUID) error {
	// TODO: re-entrance validate it isn't already confirmed

	// get auction from store
	frags, err := s.orderBookStore.GetAuction(ctx, auctionId)
	if err != nil {
		logctx.Warn(ctx, "GetAuction Failed", logger.Error(err))
		return err
	}

	orders := []*models.Order{}
	// validate all pending orders fragments of auction
	for _, frag := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		// no return during erros as what can be revert, should
		if order == nil {
			logctx.Error(ctx, "order not found while reverting an auction", logger.Error(err))
		} else if !validatePendingFrag(frag, order) {
			logctx.Error(ctx, "Auction fragments should be valid during a revert request", logger.Error(err))
		} else {
			// success
			order.SizePending.Sub(frag.Size)
			orders = append(orders, order)
		}
	}
	// store orders
	s.orderBookStore.StoreOrders(ctx, orders)

	return s.orderBookStore.RemoveAuction(ctx, auctionId)
}

func (s *Service) AuctionMined(ctx context.Context, auctionId uuid.UUID) error {
	// TODO: re-entrance validate it doesnt already confirmed

	// get auction from store
	frags, err := s.orderBookStore.GetAuction(ctx, auctionId)
	if err != nil {
		logctx.Warn(ctx, "GetAuction Failed", logger.Error(err))
		return err
	}
	var filledOrders []*models.Order

	// validate all pending orders fragments of auction
	for _, frag := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
		if order == nil {
			// cancel auction
			s.orderBookStore.RemoveAuction(ctx, auctionId) // PANIC - shouldn't happen
			logctx.Error(ctx, "Auction fragment's order should not be removed during pending to be mined", logger.Error(err))

			// return empty
			logctx.Error(ctx, err.Error())
			return models.ErrOrderNotFound
		} else if !validatePendingFrag(frag, order) {
			// cancel auction
			s.orderBookStore.RemoveAuction(ctx, auctionId) // PANIC - shouldn't happen
			logctx.Error(ctx, "Auction fragments should be valid after pending to be mined", logger.Error(err))

			logctx.Error(ctx, fmt.Sprintf("validatePendingFrag failed. PendingSize: %s FragSize:%s", order.SizePending.String(), frag.Size.String()))
			return models.ErrAuctionInvalid
		} else {
			// fill fragment in the order
			order.SizePending.Sub(frag.Size)
			order.SizeFilled.Add(frag.Size)

			// success - mark as filled
			filledOrders = append(filledOrders, order)
		}
	}

	// store orders
	// TODO: close completely filled orders
	s.orderBookStore.StoreOrders(ctx, filledOrders)

	return s.orderBookStore.RemoveAuction(ctx, auctionId) // no need to revert pending its done in line 124

}
