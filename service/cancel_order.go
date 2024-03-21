package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrderInput struct {
	Id          uuid.UUID
	IsClientOId bool
	UserId      uuid.UUID
}

// Flow chart - https://miro.com/welcomeonboard/Umt0YnpDN3BEcUh1U0JZaHNpejJNUHV3QmpBTGpTNFdybXVlemk2QlV4RHAwc2xVSXR5VzM0NzJwUlhGZEFRMnwzMDc0NDU3MzU4MzEyODA0NjQ2fDI=?share_link_id=23847173917

// CancelOrder cancels an order by its ID or clientOId. If `isClientOId` is true, the `id` is treated as a clientOinput.Id, otherwise it is treated as an orderId
func (s *Service) CancelOrder(ctx context.Context, input CancelOrderInput) (*uuid.UUID, error) {

	order, err := s.getOrder(ctx, input.IsClientOId, input.Id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		logctx.Warn(ctx, "order not found", logger.String("id", input.Id.String()), logger.Bool("isClientOId", input.IsClientOId))
		return nil, models.ErrNotFound
	}

	if order.Cancelled {
		logctx.Warn(ctx, "order already cancelled", logger.String("orderId", order.Id.String()))
		return nil, models.ErrOrderCancelled
	}

	if order.IsFilled() {
		logctx.Warn(ctx, "order already filled", logger.String("orderId", order.Id.String()))
		return nil, models.ErrOrderFilled
	}

	err = s.orderBookStore.PerformTx(ctx, func(txid uint) error {
		order.Cancelled = true

		// remove from prices
		if err := s.orderBookStore.TxModifyPrices(ctx, txid, models.Remove, *order); err != nil {
			logctx.Error(ctx, "Failed removing order from prices", logger.String("id", input.Id.String()), logger.String("side", order.Side.String()), logger.Error(err))
			return fmt.Errorf("failed removing order from prices: %w", err)
		}

		// remove from user's open orders
		if err = s.orderBookStore.TxModifyUserOpenOrders(ctx, txid, models.Remove, *order); err != nil {
			logctx.Error(ctx, "Failed removing order from user open orders", logger.String("id", input.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
			return fmt.Errorf("failed removing order from user open orders: %w", err)
		}

		switch {
		// ORDER IS PARTIALLY FILLED AND NOT PENDING
		case !order.IsUnfilled() && !order.IsPending():
			logctx.Debug(ctx, "cancelling partially filled and not pending order", logger.String("orderId", order.Id.String()))
			if err := s.orderBookStore.TxModifyUserFilledOrders(ctx, txid, models.Add, *order); err != nil {
				logctx.Error(ctx, "Failed adding order to user filled orders", logger.String("id", input.Id.String()), logger.String("userId", order.UserId.String()), logger.Error(err))
				return fmt.Errorf("failed adding order to user filled orders: %w", err)
			}
			if err = s.orderBookStore.TxModifyOrder(ctx, txid, models.Update, *order); err != nil {
				logctx.Error(ctx, "Failed updating order to cancelled", logger.String("id", input.Id.String()), logger.Error(err))
				return fmt.Errorf("failed updating order to cancelled: %w", err)
			}
		// ORDER IS PARTIALLY FILLED AND PENDING
		case !order.IsUnfilled() && order.IsPending():
			logctx.Debug(ctx, "cancelling partially filled and pending order", logger.String("orderId", order.Id.String()))
			if err = s.orderBookStore.TxModifyOrder(ctx, txid, models.Update, *order); err != nil {
				logctx.Error(ctx, "Failed updating order", logger.String("id", input.Id.String()), logger.Error(err))
				return fmt.Errorf("failed updating order: %w", err)
			}
		// ORDER IS UNFILLED AND NOT PENDING
		case order.IsUnfilled() && !order.IsPending():
			logctx.Debug(ctx, "cancelling unfilled and not pending order", logger.String("orderId", order.Id.String()))
			if err := s.orderBookStore.TxModifyClientOId(ctx, txid, models.Remove, *order); err != nil {
				logctx.Error(ctx, "Failed removing order from clientOId", logger.String("id", input.Id.String()), logger.Error(err))
				return fmt.Errorf("failed removing unfilled order: %w", err)
			}
			if err = s.orderBookStore.TxModifyOrder(ctx, txid, models.Remove, *order); err != nil {
				logctx.Error(ctx, "Failed removing order", logger.String("id", input.Id.String()), logger.Error(err))
				return fmt.Errorf("failed removing unfilled order: %w", err)
			}
		// ORDER IS UNFILLED AND PENDING
		case order.IsUnfilled() && order.IsPending():
			logctx.Debug(ctx, "cancelling unfilled and pending order", logger.String("orderId", order.Id.String()))
			if err = s.orderBookStore.TxModifyOrder(ctx, txid, models.Update, *order); err != nil {
				logctx.Error(ctx, "Failed updating order", logger.String("id", input.Id.String()), logger.Error(err))
				return fmt.Errorf("failed updating order: %w", err)
			}
		default:
			logctx.Error(ctx, "unexpected order state", logger.String("orderId", order.Id.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))
			return models.ErrUnexpectedError
		}

		return nil
	})

	logctx.Debug(ctx, "order cancelled", logger.String("orderId", order.Id.String()), logger.String("userId", order.UserId.String()), logger.String("size", order.Size.String()), logger.String("sizeFilled", order.SizeFilled.String()), logger.String("sizePending", order.SizePending.String()))

	s.publishOrderEvent(ctx, order)

	return &order.Id, nil
}

func (s *Service) getOrder(ctx context.Context, isClientOId bool, orderId uuid.UUID) (order *models.Order, err error) {
	if isClientOId {
		logctx.Debug(ctx, "finding order by clientOId", logger.String("clientOId", orderId.String()))

		order, err = s.orderBookStore.FindOrderById(ctx, orderId, true)
		if err != nil {
			logctx.Error(ctx, "could not get order ID by clientOId", logger.Error(err))
			return nil, err
		}
	} else {
		logctx.Debug(ctx, "finding order by orderId", logger.String("orderId", orderId.String()))
		order, err = s.orderBookStore.FindOrderById(ctx, orderId, false)
		if err != nil {
			logctx.Error(ctx, "could not get order by orderId", logger.Error(err))
			return nil, err
		}
	}

	return order, nil
}
