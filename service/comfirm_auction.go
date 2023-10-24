package service

import (
	"context"
)

// orderID->amount bought or sold in A token always

type ConfirmOrderRes struct {
	Orders []string
}

func (s *Service) ConfirmAuction(ctx context.Context, auctionId string) (ConfirmOrderRes, error) {
	return ConfirmOrderRes{}, nil
}
