package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type ProcessOrderInput struct {
	UserId        uuid.UUID
	Price         decimal.Decimal
	Symbol        models.Symbol
	Size          decimal.Decimal
	Side          models.Side
	ClientOrderID uuid.UUID
}

var (
	ErrClashingOrderId = errors.New("order with that ID already exists")
)

func (s *Service) ProcessOrder(ctx context.Context, input ProcessOrderInput) (models.Order, error) {

	// TODO: logic for when order can be immediately filled and does not need to be added to the order book

	existingOrder, err := s.orderBookStore.FindOrderById(ctx, input.ClientOrderID, true)

	if err != nil && err != models.ErrOrderNotFound {
		logctx.Error(ctx, "unexpected error when finding order by clientOrderId", logger.Error(err))
		return models.Order{}, models.ErrUnexpectedError
	}

	if existingOrder == nil {
		logctx.Info(ctx, "no existing order with same orderId. Trying to create new order", logger.String("clientOrderId", input.ClientOrderID.String()))
		return s.createNewOrder(ctx, input, input.ClientOrderID)
	}

	if existingOrder.UserId != input.UserId {
		logctx.Warn(ctx, "order already exists with different userId", logger.Error(err))
		return models.Order{}, ErrClashingOrderId
	}

	if existingOrder.ClientOId == input.ClientOrderID {
		logctx.Warn(ctx, "order already exists with same clientOrderId", logger.Error(err), logger.String("clientOrderId", input.ClientOrderID.String()))
		return models.Order{}, models.ErrOrderAlreadyExists
	}

	logctx.Error(ctx, "did not follow any cases when processing order", logger.String("clientOrderId", input.ClientOrderID.String()), logger.String("userId", input.UserId.String()), logger.String("price", input.Price.String()), logger.String("size", input.Size.String()), logger.String("symbol", input.Symbol.String()), logger.String("side", input.Side.String()))

	return models.Order{}, models.ErrUnexpectedError
}

func (s *Service) createNewOrder(ctx context.Context, input ProcessOrderInput, clientOId uuid.UUID) (models.Order, error) {
	orderId := uuid.New()

	logctx.Info(ctx, "creating new order", logger.String("orderId", orderId.String()), logger.String("clientOrderId", clientOId.String()))

	order := models.Order{
		Id:        orderId,
		ClientOId: clientOId,
		UserId:    input.UserId,
		Price:     input.Price,
		Symbol:    input.Symbol,
		Size:      input.Size,
		Signature: "",
		Status:    models.STATUS_OPEN,
		Side:      input.Side,
		Timestamp: time.Now().UTC(),
	}

	if err := s.orderBookStore.StoreOrder(ctx, order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "new order created", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))
	return order, nil
}
