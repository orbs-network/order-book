package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrder(ctx context.Context, orderId uuid.UUID) error {

	user := utils.GetUser(ctx)
	if user == nil {
		logctx.Error(ctx, "user not found in context")
		return models.ErrNoUserInContext
	}

	order, err := s.orderBookStore.FindOrderById(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "error occured when finding order", logger.Error(err))
		return err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("orderId", orderId.String()))
		return models.ErrOrderNotFound
	}

	if order.Status != models.STATUS_OPEN {
		logctx.Warn(ctx, "trying to cancel order that is not open", logger.String("orderId", orderId.String()), logger.String("status", order.Status.String()))
		return models.ErrOrderNotOpen
	}

	if user.ID != order.UserId {
		logctx.Warn(ctx, "user trying to cancel another user's order", logger.String("orderId", orderId.String()), logger.String("requestUserId", user.ID.String()), logger.String("orderUserId", order.UserId.String()))
		return models.ErrUnauthorized
	}

	err = s.orderBookStore.RemoveOrder(ctx, *order)
	if err != nil {
		logctx.Error(ctx, "error occured when removing order", logger.Error(err))
		return err
	}

	logctx.Info(ctx, "order removed", logger.String("orderId", orderId.String()))
	return nil
}
