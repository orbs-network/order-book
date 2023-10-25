package service

import (
	"context"
)

// orderID->amount bought or sold in A token always

type ConfirmAuctionRes struct {
	Orders []string
}

func (s *Service) ConfirmAuction(ctx context.Context, auctionId string) (ConfirmAuctionRes, error) {
	// get auction from store

	// validate all orders of auction

	// return orders signatures

	// add oredebook signature on the buffer

	return ConfirmAuctionRes{}, nil
}

func (s *Service) RemoveAuction(ctx context.Context, auctionId string) error {
	return nil
}
