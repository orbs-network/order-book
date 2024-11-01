package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) GetOpenOrders(ctx context.Context, userId uuid.UUID, symbol models.Symbol) (orders []models.Order, totalOrders int, err error) {
	logctx.Debug(ctx, "getting open orders for user", logger.String("user_id", userId.String()))
	orders, _, err = s.orderBookStore.GetOpenOrders(ctx, userId, symbol)

	if err != nil {
		logctx.Error(ctx, "error getting open orders for user", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, 0, fmt.Errorf("error getting open orders for user: %w", err)
	}

	// filter out cancelled or fully filled orders and keep only open
	res := []models.Order{}
	for _, o := range orders {
		if !o.Cancelled && !o.IsFilled() {
			res = append(res, o)
		}
	}

	logctx.Debug(ctx, "returning open orders for user", logger.String("user_id", userId.String()), logger.Int("orders_count", len(orders)))

	return res, len(res), nil
}

const MAX_FILLS = 256

// get all fills from all swaps of the user in a given time range
func (s *Service) GetSwapFills(ctx context.Context, userId uuid.UUID, symbol models.Symbol, startAt, endAt time.Time) ([]models.Fill, error) {
	logctx.Debug(ctx, "getting open orders for user", logger.String("user_id", userId.String()))

	swapIds, err := s.orderBookStore.GetUserResolvedSwapIds(ctx, userId)
	if err != nil {
		logctx.Error(ctx, "error getting user resolve swapIds", logger.Error(err), logger.String("user_id", userId.String()))
		return nil, err
	}

	if len(swapIds) == 0 {
		logctx.Warn(ctx, "user has no resolved swaps", logger.Error(err), logger.String("user_id", userId.String()))
		return []models.Fill{}, nil
	}

	fills := []models.Fill{}

	// fetch swaps
	for _, id := range swapIds {
		uid, err := uuid.Parse(id)
		if err != nil {
			logctx.Error(ctx, "failed to parse swapID", logger.Error(err), logger.String("user_id", userId.String()), logger.String("swap_id", id))
			return nil, err
		}
		// get resolved swaps
		swap, err := s.orderBookStore.GetSwap(ctx, uid, false)
		if err != nil {
			logctx.Error(ctx, "error getting a swap", logger.Error(err), logger.String("user_id", userId.String()), logger.String("swap_id", id))
			return nil, err
		}
		// check if swap is in time range
		if (swap.Resolved.Equal(startAt) || swap.Resolved.After(startAt)) && swap.Resolved.Before(endAt) {
			// iterate through fragments
			for _, frag := range swap.Frags {
				order, err := s.orderBookStore.FindOrderById(ctx, frag.OrderId, false)
				if err != nil {
					logctx.Warn(ctx, "error getting a order", logger.Error(err), logger.String("user_id", userId.String()), logger.String("order_id", id))
				}
				fills = append(fills, *models.NewFill(symbol, *swap, frag, order))

				if len(fills) >= MAX_FILLS {
					return nil, models.ErrMaxRecExceeded
				}
			}
		}
	}

	return fills, nil
}
