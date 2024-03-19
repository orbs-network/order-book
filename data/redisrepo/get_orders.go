package redisrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *redisRepository) GetOpenOrderIds(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	userOrdersKey := CreateUserOpenOrdersKey(userId)

	// Fetch all order IDs for the user
	orderIdStrs, err := r.client.ZRange(ctx, userOrdersKey, 0, -1).Result()
	if err != nil {
		logctx.Error(ctx, "failed to get order IDs for user", logger.String("userId", userId.String()), logger.Error(err))
		return nil, fmt.Errorf("failed to get order IDs for user: %v", err)
	}

	if len(orderIdStrs) == 0 {
		logctx.Warn(ctx, "no orders found for user", logger.String("userId", userId.String()))
		return nil, models.ErrNotFound
	}
	// Convert string IDs to UUIDs
	var orderIds []uuid.UUID
	for _, orderIdStr := range orderIdStrs {
		orderId, err := uuid.Parse(orderIdStr)
		if err != nil {
			logctx.Error(ctx, "failed to parse order ID", logger.String("orderId", orderIdStr), logger.Error(err))
			return nil, fmt.Errorf("failed to parse order ID: %v", err)
		}
		orderIds = append(orderIds, orderId)
	}
	return orderIds, nil

}

// func (r *redisRepository) GetOpenOrders(ctx context.Context, userId uuid.UUID, symbol models.Symbol) ([]models.Order, error) {

// 	orderIds := r.GetOpenOrdersIds()
// 	// We only want to fetch open orders
// 	orders, err := r.FindOrdersByIds(ctx, orderIds, true)
// 	if err != nil {
// 		logctx.Error(ctx, "failed to find orders by IDs", logger.String("userId", userId.String()), logger.Error(err))
// 		return nil, fmt.Errorf("failed to find orders by IDs: %v", err)
// 	}
// 	// no symbol filter
// 	if symbol == "" {
// 		return orders, nil
// 	}
// 	// filter only relevant symbols
// 	res := []models.Order{}
// 	for _, order := range orders {
// 		if order.Symbol == symbol {
// 			res = append(res, order)
// 		}
// 	}
// 	return res, nil

// }
