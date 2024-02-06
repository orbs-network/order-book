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

// CancelOrder cancels an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOinput.Id, otherwise it is treated as an orderId.
func (s *Service) CancelOrder(ctx context.Context, input CancelOrderInput) (cancelledOrderId *uuid.UUID, err error) {

	order, err := s.getOrder(ctx, input.IsClientOId, input.Id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, models.ErrNotFound
	}

	if order.IsPending() {
		logctx.Warn(ctx, "cancelling order not possible when order is pending", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		return nil, models.ErrOrderPending
	}

	if order.IsFilled() {
		logctx.Warn(ctx, "cancelling order not possible when order is filled", logger.String("orderId", order.Id.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("size", order.Size.String()))
		return nil, models.ErrOrderFilled
	}

	if order.IsUnfilled() {
		err = s.orderBookStore.CancelUnfilledOrder(ctx, *order)
		if err != nil {
			logctx.Error(ctx, "error occured when removing order", logger.Error(err))
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
