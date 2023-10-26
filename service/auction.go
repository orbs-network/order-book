package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// orderID->amount bought or sold in A token always

type ConfirmAuctionRes struct {
	Orders []string
}

func (s *Service) validateFillOrder(fillOrder string) *models.Order {
	// get order from ID

	// check if order is still open
	// order.size - filledOrder.filled >= fillOrder.amount
	return nil
}

func (s *Service) ConfirmAuction(ctx context.Context, auctionId string) (ConfirmAuctionRes, error) {
	// get auction from store
	arrOfFillOrders := s.orderBookStore.GetAuction(auctionId)

	// validate all orders of auction
	for fillOrder := range arrOfFillOrders {
		// get order by ID
		// validate order
		// if order is valid
		//order := s.validateFillOrder(models.Order{})

		// validate
		// if order == nil {
		// 	logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		// 	return ConfirmAuctionRes{}, models.ErrInsufficientLiquity
		// }

	}

	// return orders signatures

	// add oredebook signature on the buffer

	// lock funds

	// remove auction store
	//s.orderBookStore.RemoveAuction(auctionId)

	// error
	logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
	return ConfirmAuctionRes{}, models.ErrInsufficientLiquity

	return ConfirmAuctionRes{}, nil
}

func (s *Service) RemoveAuction(ctx context.Context, auctionId string) error {
	return nil
}
