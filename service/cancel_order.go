package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrderInput struct {
	Id          uuid.UUID
	IsClientOId bool
	UserId      uuid.UUID
}

//flow chart
//https://miro.com/welcomeonboard/Umt0YnpDN3BEcUh1U0JZaHNpejJNUHV3QmpBTGpTNFdybXVlemk2QlV4RHAwc2xVSXR5VzM0NzJwUlhGZEFRMnwzMDc0NDU3MzU4MzEyODA0NjQ2fDI=?share_link_id=23847173917

func (s *Service) CancelOrder(ctx context.Context, input CancelOrderInput) (*uuid.UUID, error) {
	// get order
	order, err := s.getOrder(ctx, input.IsClientOId, input.Id)
	if err != nil {
		logctx.Warn(ctx, "order not found", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, err
	}
	// remove from price
	txid, err := s.orderBookStore.TxStart(ctx)
	if err != nil {
		return nil, err
	}
	defer s.orderBookStore.TxEnd(ctx, txid)

	err = s.orderBookStore.TxRemoveOrderFromPrice(ctx, txid, *order)
	if err != nil {
		logctx.Error(ctx, "TxRemoveOrderFromPrice", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, err
	}
	// mark as cancelled
	order.Cancelled = true
	// isUnfilled?
	if order.IsUnfilled() {
		// non Pending?
		if !order.IsPending() {
			err := s.orderBookStore.TxDeleteOrder(ctx, txid, order.Id)
			if err != nil {
				logctx.Error(ctx, "DeleteOrder error  in Cancel Unfilled non pending Order", logger.String("orderId", input.Id.String()), logger.Error(err))
				return nil, err
			}
			return &order.Id, nil
		}
	}
	// store cancelled state change
	err = s.orderBookStore.TxStoreOrder(ctx, txid, *order)
	if err != nil {
		logctx.Error(ctx, "StoreOpenOrder error  in cancel of filled or pending order", logger.String("orderId", input.Id.String()), logger.Error(err))
		return nil, err
	}
	return &order.Id, nil
}

// CancelOrder cancels an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOinput.Id, otherwise it is treated as an orderId.
func (s *Service) CancelOrderOld(ctx context.Context, input CancelOrderInput) (*uuid.UUID, error) {

	order, err := s.getOrder(ctx, input.IsClientOId, input.Id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, models.ErrNotFound
	}

	if order.IsPending() {
		logctx.Info(ctx, "cancelling a pending order", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		err = s.orderBookStore.CancelPendingOrder(ctx, *order)
		if err != nil {
			logctx.Error(ctx, "error CancelPendingOrder", logger.Error(err))
			return nil, err
		}
		return nil, models.ErrOrderPending
	}

	if order.IsFilled() {
		logctx.Warn(ctx, "cancelling order not possible when order is filled", logger.String("orderId", order.Id.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("size", order.Size.String()))
		return nil, models.ErrOrderFilled
	}

	if order.IsUnfilled() {
		err = s.orderBookStore.CancelUnfilledOrder(ctx, *order)
		if err != nil {
			logctx.Error(ctx, "error CancelUnfilledOrder", logger.Error(err))
			return nil, err
		}

		logctx.Info(ctx, "unfilled order removed", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))

		return &order.Id, nil
	} else {
		err = s.orderBookStore.CancelPartialFilledOrder(ctx, *order)
		if err != nil {
			logctx.Error(ctx, "error occured when cancelling partial order", logger.Error(err))
			return nil, err
		}

		logctx.Info(ctx, "partial filled order cancelled", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))

		return &order.Id, nil
	}
}

func (s *Service) getOrder(ctx context.Context, isClientOId bool, orderId uuid.UUID) (order *models.Order, err error) {
	if isClientOId {
		logctx.Info(ctx, "finding order by clientOId", logger.String("clientOId", orderId.String()))

		order, err = s.orderBookStore.FindOrderById(ctx, orderId, true)
		if err != nil {
			logctx.Error(ctx, "could not get order ID by clientOId", logger.Error(err))
			return nil, err
		}
	} else {
		logctx.Info(ctx, "finding order by orderId", logger.String("orderId", orderId.String()))
		order, err = s.orderBookStore.FindOrderById(ctx, orderId, false)
		if err != nil {
			logctx.Error(ctx, "could not get order by orderId", logger.Error(err))
			return nil, err
		}
	}

	return order, nil
}
