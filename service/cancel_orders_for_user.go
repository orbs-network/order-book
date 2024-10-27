package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) CancelOrdersForUser(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orderIds []uuid.UUID, err error) {

	orders, err := s.orderBookStore.GetOpenOrdersForUser(ctx, userId)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, err
		}
		logctx.Error(ctx, "could not get open orders for user", logger.Error(err), logger.String("userId", userId.String()))
		return nil, err
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
			// error
			if err != nil {
				logctx.Error(ctx, "could not cancel order", logger.Error(err), logger.String("orderId", order.Id.String()))
			} else if uid != nil {
				// success
				res = append(res, *uid)
			}

		}
	}

	return res, nil
}
