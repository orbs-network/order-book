package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) AddOrder(ctx context.Context, userId uuid.UUID, price decimal.Decimal, symbol models.Symbol, size decimal.Decimal) (models.Order, error) {

	existingOrder, err := s.orderBookStore.FindOrder(ctx, models.FindOrderInput{
		UserId: userId,
		Price:  price,
		Symbol: symbol,
	})
	if err != nil {
		logctx.Error(ctx, "error occured when finding order", logger.Error(err))
		return models.Order{}, err
	}

	if existingOrder != nil {
		logctx.Warn(ctx, "user already has an order at this price", logger.String("userID", userId.String()), logger.String("price", price.String()))
		return models.Order{}, models.ErrOrderAlreadyExists

	}

	id := uuid.New()
	order := models.Order{
		Id:        id,
		UserId:    userId,
		Price:     price,
		Symbol:    symbol,
		Size:      size,
		Signature: nil,
		Status:    models.STATUS_OPEN,
	}

	if err := s.orderBookStore.StoreOrder(ctx, order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "order added", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))
	return order, nil
}

func (s *Service) CancelOrder(ctx context.Context, orderId uuid.UUID) error {

	order, err := s.orderBookStore.FindOrderById(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "error occured when finding order", logger.Error(err))
		return err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("orderId", orderId.String()))
		return models.ErrOrderNotFound
	}

	err = s.orderBookStore.RemoveOrder(ctx, orderId)
	if err != nil {
		logctx.Error(ctx, "error occured when removing order", logger.Error(err))
		return err
	}

	logctx.Info(ctx, "order removed", logger.String("orderId", orderId.String()))
	return nil
}
