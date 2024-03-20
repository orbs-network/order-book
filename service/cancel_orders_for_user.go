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
		logctx.Debug(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return []uuid.UUID{}, err
	}

	orders, err := s.orderBookStore.FindOrdersByIds(ctx, ids, true)
	if err != nil {
		logctx.Error(ctx, "FindOrdersByIds open orders failed", logger.Error(err), logger.String("userId", userId.String()))
		return []uuid.UUID{}, err
	}
	res := []uuid.UUID{}
	for _, order := range orders {
		// matching symbol only if provided
		if symbol == "" || order.Symbol == symbol {
			uid, err := s.CancelOrder(ctx, CancelOrderInput{
				Id:          order.Id,
				IsClientOId: false,
				UserId:      userId,
			})
			if err != nil {
				logctx.Error(ctx, "could not cancel order", logger.Error(err), logger.String("orderId", uid.String()))
			}
			res = append(res, *uid)
		}
	}

	return res, nil
}
