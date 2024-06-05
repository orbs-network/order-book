package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/abi"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type CreateOrderInput struct {
	UserId        uuid.UUID
	Price         decimal.Decimal
	Symbol        models.Symbol
	Size          decimal.Decimal
	Side          models.Side
	ClientOrderID uuid.UUID
	Eip712Sig     string
	AbiFragment   abi.Order
}

func (s *Service) CreateOrder(ctx context.Context, input CreateOrderInput) (models.Order, error) {

	existingOrder, err := s.orderBookStore.FindOrderById(ctx, input.ClientOrderID, true)

	if err != nil && err != models.ErrNotFound {
		logctx.Error(ctx, "unexpected error when finding order by clientOrderId", logger.Error(err))
		return models.Order{}, fmt.Errorf("unexpected error when finding order by clientOrderId: %s", err)
	}

	if existingOrder == nil {
		return s.createNewOrder(ctx, input, input.UserId)
	}

	if existingOrder.UserId != input.UserId {
		logctx.Warn(ctx, "order already exists with different userId", logger.Error(err))
		return models.Order{}, models.ErrClashingOrderId
	}

	if existingOrder.ClientOId == input.ClientOrderID {
		logctx.Warn(ctx, "order already exists with same clientOrderId", logger.Error(err), logger.String("clientOrderId", input.ClientOrderID.String()))
		return models.Order{}, models.ErrClashingClientOrderId
	}

	logctx.Error(ctx, "did not follow any cases when creating order", logger.String("clientOrderId", input.ClientOrderID.String()), logger.String("userId", input.UserId.String()), logger.String("price", input.Price.String()), logger.String("size", input.Size.String()), logger.String("symbol", input.Symbol.String()), logger.String("side", input.Side.String()))

	return models.Order{}, models.ErrUnexpectedError
}

func (s *Service) createNewOrder(ctx context.Context, input CreateOrderInput, userId uuid.UUID) (models.Order, error) {
	orderId := uuid.New()

	logctx.Debug(ctx, "creating new order", logger.String("orderId", orderId.String()), logger.String("clientOrderId", input.ClientOrderID.String()))

	// validate price
	if input.Price.IsZero() || input.Price.IsNegative() {
		logctx.Warn(ctx, "price has to be positive", logger.String("orderId", orderId.String()), logger.String("price", input.Price.String()))
		return models.Order{}, models.ErrInvalidInput
	}
	// validate cross trade
	depth, err := s.GetMarketDepth(ctx, input.Symbol, 1)
	if err != nil {
		logctx.Warn(ctx, "market depth failed", logger.String("orderId", orderId.String()), logger.String("price", input.Price.String()))
	}
	if len(depth.Asks) > 0 {
		minAsk := depth.Asks[0][0]
		if input.Side == models.BUY && input.Price.GreaterThanOrEqual(minAsk) {
			logctx.Warn(ctx, "CrossTrade order rejected. bid price is higher than minAsk", logger.String("orderId", orderId.String()), logger.String("price", input.Price.String()), logger.String("min_ask", minAsk.String()))
			return models.Order{}, models.ErrCrossTrade
		}
	}

	szBids := len(depth.Bids)
	if szBids > 0 {
		maxBid := depth.Bids[0][szBids-1]
		if input.Side == models.SELL && input.Price.LessThanOrEqual(maxBid) {
			logctx.Warn(ctx, "CrossTrade order rejected. bid price is higher than minAsk", logger.String("orderId", orderId.String()), logger.String("price", input.Price.String()), logger.String("maxBid", maxBid.String()))
			return models.Order{}, models.ErrCrossTrade
		}
	}

	// validate size
	if input.Size.IsZero() || input.Size.IsNegative() {
		logctx.Warn(ctx, "size has to be positive", logger.String("orderId", orderId.String()), logger.String("size", input.Size.String()))
		return models.Order{}, models.ErrInvalidInput
	}

	order := models.Order{
		Id:        orderId,
		ClientOId: input.ClientOrderID,
		UserId:    userId,
		Price:     input.Price,
		Symbol:    input.Symbol,
		Size:      input.Size,
		Signature: models.Signature{
			Eip712Sig:   input.Eip712Sig,
			AbiFragment: input.AbiFragment,
		},
		Side:      input.Side,
		Timestamp: time.Now().UTC(),
	}

	if err := s.orderBookStore.StoreOpenOrder(ctx, order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	s.publishOrderEvent(ctx, &order)

	return order, nil
}
