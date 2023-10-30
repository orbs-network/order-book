package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// CancelOrder cancels an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOId, otherwise it is treated as an orderId.
func (s *Service) CancelOrder(ctx context.Context, id uuid.UUID, isClientOId bool) (cancelledOrderId *uuid.UUID, err error) {

	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user not found in context")
		return nil, models.ErrNoUserInContext
	}

	var order *models.Order

	if isClientOId {
		logctx.Info(ctx, "finding order by clientOId", logger.String("clientOId", id.String()))

		order, err = s.orderBookStore.FindOrderById(ctx, id, true)
		if err != nil {
			logctx.Error(ctx, "could not get order ID by clientOId", logger.Error(err))
			return nil, err
		}
	} else {
		logctx.Info(ctx, "finding order by orderId", logger.String("orderId", id.String()))
		order, err = s.orderBookStore.FindOrderById(ctx, id, false)
		if err != nil {
			logctx.Error(ctx, "could not get order by orderId", logger.Error(err))
			return nil, err
		}
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("id", id.String()), logger.Bool("isClientOId", isClientOId))
		return nil, models.ErrOrderNotFound
	}

	if order.Status != models.STATUS_OPEN {
		logctx.Warn(ctx, "trying to cancel order that is not open", logger.String("orderId", order.Id.String()), logger.String("status", order.Status.String()))
		return nil, models.ErrOrderNotOpen
	}

	if user.ID != order.UserId {
		logctx.Warn(ctx, "user trying to cancel another user's order", logger.String("orderId", order.Id.String()), logger.String("requestUserId", user.ID.String()), logger.String("orderUserId", order.UserId.String()))
		return nil, models.ErrUnauthorized
	}

	err = s.orderBookStore.RemoveOrder(ctx, *order)
	if err != nil {
		logctx.Error(ctx, "error occured when removing order", logger.Error(err))
		return nil, err
	}

	logctx.Info(ctx, "order removed", logger.String("orderId", order.Id.String()))
	return &order.Id, nil
}
