package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrdersForUser(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orderIds []uuid.UUID, err error) {

	ids, err := s.orderBookStore.GetOpenOrderIds(ctx, userId)
	if err != nil {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return []uuid.UUID{}, err
	}
	res := []uuid.UUID{}
	for _, id := range ids {
		order, err := s.getOrder(ctx, false, id)
		if err != nil {
			logctx.Error(ctx, "could not get order", logger.Error(err), logger.String("orderId", id.String()))
		}
		if order != nil {
			// matching symbol only if provided
			if symbol == "" || order.Symbol == symbol {
				uid, err := s.CancelOrder(ctx, CancelOrderInput{
					Id:          id,
					IsClientOId: false,
					UserId:      userId,
				}) //, symbol)
				if err != nil {
					logctx.Error(ctx, "could not cancel order", logger.Error(err), logger.String("orderId", uid.String()))
				}
				res = append(res, *uid)
			}
		}
	}

	return res, nil
}
