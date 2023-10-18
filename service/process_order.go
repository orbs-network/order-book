package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type ProcessOrderInput struct {
	UserId uuid.UUID
	Price  decimal.Decimal
	Symbol models.Symbol
	Size   decimal.Decimal
	Side   models.Side
	// Optional user provided order ID. We still generate a seperate ID for the order
	ClientOrderID *uuid.UUID
}

func (s *Service) ProcessOrder(ctx context.Context, input ProcessOrderInput) (models.Order, error) {

	// TODO: add to existing order if possible

	id := uuid.New()
	order := models.Order{
		Id:            id,
		UserId:        input.UserId,
		Price:         input.Price,
		Symbol:        input.Symbol,
		Size:          input.Size,
		Signature:     "",
		Status:        models.STATUS_OPEN,
		Side:          input.Side,
		ClientOrderID: generateClientOrderId(input.ClientOrderID, id),
	}

	if err := s.orderBookStore.StoreOrder(ctx, order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "order added", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))
	return order, nil
}

func generateClientOrderId(clientOrderId *uuid.UUID, orderId uuid.UUID) uuid.UUID {
	if clientOrderId == nil {
		return orderId
	}

	return *clientOrderId
}
