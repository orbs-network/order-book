package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/featureflags"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
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
	Eip712MsgData map[string]interface{}
}

func (s *Service) CreateOrder(ctx context.Context, input CreateOrderInput) (models.Order, error) {

	if featureflags.ShouldVerifySig == "" || featureflags.ShouldVerifySig == "true" {
		logctx.Info(ctx, "verifying signature", logger.String("userId", input.UserId.String()))

		user := utils.GetUserCtx(ctx)
		if user == nil {
			logctx.Error(ctx, "user should be in context")
			return models.Order{}, fmt.Errorf("user should be in context")
		}

		isVerifed, err := s.blockchainClient.VerifySignature(ctx, VerifySignatureInput{
			MessageData: input.Eip712MsgData,
			Signature:   input.Eip712Sig,
			PublicKey:   user.PubKey,
		})

		if err != nil {
			logctx.Warn(ctx, "signature verification error", logger.Error(err), logger.String("userId", user.Id.String()))
			return models.Order{}, models.ErrSignatureVerificationError
		}

		if !isVerifed {
			logctx.Warn(ctx, "signature verification failed", logger.String("userId", user.Id.String()))
			return models.Order{}, models.ErrSignatureVerificationFailed
		}
	}

	existingOrder, err := s.orderBookStore.FindOrderById(ctx, input.ClientOrderID, true)

	if err != nil && err != models.ErrNotFound {
		logctx.Error(ctx, "unexpected error when finding order by clientOrderId", logger.Error(err))
		return models.Order{}, fmt.Errorf("unexpected error when finding order by clientOrderId: %s", err)
	}

	if existingOrder == nil {
		logctx.Info(ctx, "no existing order with same orderId. Trying to create new order", logger.String("clientOrderId", input.ClientOrderID.String()))
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

	logctx.Info(ctx, "creating new order", logger.String("orderId", orderId.String()), logger.String("clientOrderId", input.ClientOrderID.String()))

	order := models.Order{
		Id:        orderId,
		ClientOId: input.ClientOrderID,
		UserId:    userId,
		Price:     input.Price,
		Symbol:    input.Symbol,
		Size:      input.Size,
		Signature: models.Signature{
			Eip712Sig:     input.Eip712Sig,
			Eip712MsgData: input.Eip712MsgData,
		},
		Side:      input.Side,
		Timestamp: time.Now().UTC(),
	}

	if err := s.orderBookStore.StoreOpenOrder(ctx, order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "new order created", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))

	s.publishOrderEvent(ctx, &order, models.STATUS_OPEN)

	return order, nil
}
