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
	FillReqs      []*models.OrderFrag
	BookSignature []byte
}

func validateFillReq(fillReq models.OrderFrag, order *models.Order) bool {

	// check if order is still open
	if order.Status != models.STATUS_OPEN {
		return false
	}
	// order.size - (Order.filled + prder.pending) >= fillOrder.amount
	orderLockedSum := order.SizeFilled.Sub(order.SizePending)
	return order.Size.Sub(orderLockedSum).GreaterThanOrEqual(fillReq.Amount)
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
	for _, fillReq := range frags {
		// get order by ID
		order, err := s.orderBookStore.FindOrderById(ctx, fillReq.OrderId, false)
		if order == nil {
			// cancel auction
			s.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrOrderNotFound
		} else if !validateFillReq(fillReq, order) {
			// cancel auction
			s.RemoveAuction(ctx, auctionId)

			// return empty
			logctx.Warn(ctx, err.Error())
			return ConfirmAuctionRes{}, models.ErrAuctionInvalid
		} else {
			// success- append
			res.Orders = append(res.Orders, order)
			res.FillReqs = append(res.FillReqs, &fillReq)

			// later s.orderBookStore.FillOrder()

		}
	}
	// process all fill requests
	for i := 0; i < len(res.Orders); i++ {
		// lock fillReq.Amount as pending per order - no STATUS_PENDING is needed
		//s.orderBookStore.SetPendingOrders(ctx, res.Orders[i], res.FillReqs[i].Amount)
		res.Orders[i].SizePending = res.FillReqs[i].Amount

		// s.ProcessOrder(ctx, ProcessOrderInput{
		// 	UserId:        uuid.UUID{},
		// 	Price:         res.Orders[i].Price,
		// 	Size:          res.FillReqs[i].Amount,
		// 	Side:          res.Orders[i].Side,
		// 	ClientOrderID: res.Orders[i].ClientOId,
		// })
	}
	s.orderBookStore.StoreOrders(ctx, res.Orders)

	// type ProcessOrderInput struct {
	// 	UserId        uuid.UUID
	// 	Price         decimal.Decimal
	// 	Symbol        models.Symbol
	// 	Size          decimal.Decimal
	// 	Side          models.Side
	// 	ClientOrderID uuid.UUID
	// }

	// type Order struct {
	// 	Id        uuid.UUID       `json:"orderId"`
	// 	ClientOId uuid.UUID       `json:"clientOrderId"`
	// 	UserId    uuid.UUID       `json:"userId"`
	// 	Price     decimal.Decimal `json:"price"`
	// 	Symbol    Symbol          `json:"symbol"`
	// 	Size      decimal.Decimal `json:"size"`
	// 	Signature string          `json:"-" ` // EIP 712
	// 	Status    Status          `json:"-"`  // when order is pending, it should not be updateable
	// 	Side      Side            `json:"side"`
	// 	Timestamp time.Time       `json:"timestamp"`
	// }

	// return orders signatures

	// add oredebook signature on the buffer
	res.BookSignature = []byte("todo:sign")

	// lock funds
	//s.orderBookStore.

	// set entire auction as pending ??
	//s.orderBookStore.RemoveAuction(auctionId)

	// error

	return res, nil
}

func (s *Service) RemoveAuction(ctx context.Context, auctionId uuid.UUID) error {
	return nil
}
