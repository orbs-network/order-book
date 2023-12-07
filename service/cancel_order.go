package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type CancelOrderInput struct {
	Id          uuid.UUID
	IsClientOId bool
	UserId      uuid.UUID
}

// CancelOrder cancels an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOinput.Id, otherwise it is treated as an orderId.
func (s *Service) CancelOrder(ctx context.Context, input CancelOrderInput) (cancelledOrderId *uuid.UUID, err error) {

	var order *models.Order

	if input.IsClientOId {
		logctx.Info(ctx, "finding order by clientOId", logger.String("clientOId", input.Id.String()))

		order, err = s.orderBookStore.FindOrderById(ctx, input.Id, true)
		if err != nil {
			logctx.Error(ctx, "could not get order ID by clientOId", logger.Error(err))
			return nil, err
		}
	} else {
		logctx.Info(ctx, "finding order by orderId", logger.String("orderId", input.Id.String()))
		order, err = s.orderBookStore.FindOrderById(ctx, input.Id, false)
		if err != nil {
			logctx.Error(ctx, "could not get order by orderId", logger.Error(err))
			return nil, err
		}
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, models.ErrNotFound
	}

	if order.SizePending.GreaterThan(decimal.Zero) {
		logctx.Warn(ctx, "cancelling order not possible when order is pending", logger.String("orderId", order.Id.String()), logger.String("sizePending", order.SizePending.String()))
		return nil, models.ErrOrderPending
	}

	if order.Status != models.STATUS_OPEN {
		logctx.Warn(ctx, "trying to cancel order that is not open", logger.String("orderId", order.Id.String()), logger.String("status", order.Status.String()))
		return nil, models.ErrOrderNotOpen
	}

	if input.UserId != order.UserId {
		logctx.Warn(ctx, "user trying to cancel another user's order", logger.String("orderId", order.Id.String()), logger.String("requestUserId", input.UserId.String()), logger.String("orderUserId", order.UserId.String()))
		return nil, models.ErrUnauthorized
	}

	err = s.orderBookStore.RemoveOrder(ctx, *order)
	if err != nil {
		logctx.Error(ctx, "error occured when removing order", logger.Error(err))
		return nil, err
	}

	logctx.Info(ctx, "order removed", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))
	return &order.Id, nil
}
