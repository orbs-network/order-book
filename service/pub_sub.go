package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) SubscribeUserOrders(ctx context.Context, userId uuid.UUID) (chan []byte, error) {
	logctx.Debug(ctx, "subscribing to user orders", logger.String("userId", userId.String()))

	eventKey := models.CreateUserOrdersEventKey(userId)

	channel, err := s.orderBookStore.SubscribeToEvents(ctx, fmt.Sprintf("user_orders:%s", userId))
	if err != nil {
		logctx.Error(ctx, "failed to subscribe to user orders", logger.String("event", eventKey), logger.Error(err))
		return nil, fmt.Errorf("failed to subscribe to user orders: %w", err)
	}

	return channel, nil
}

func publishFillEvent(ctx context.Context, store store.OrderBookStore, userId uuid.UUID, fill models.Fill) {
	//value, err := json.Marshal(fill)
	value, err := json.Marshal(struct {
		Event string `json:"event"`
		models.Fill
	}{
		Event: "order-fill",
		Fill:  fill,
	})
	if err != nil {
		logctx.Error(ctx, "failed to marshal order to json", logger.Error(err))
	}

	key := models.CreateUserOrdersEventKey(userId)
	if err := store.PublishEvent(ctx, key, value); err != nil {
		logctx.Error(ctx, "failed to publish fill event", logger.String("event", key), logger.Error(err))
	}

}
func (e *EvmClient) publishFillEvent(ctx context.Context, userId uuid.UUID, fill models.Fill) {
	publishFillEvent(ctx, e.orderBookStore, userId, fill)
}

func (s *Service) publishFillEvent(ctx context.Context, userId uuid.UUID, fill models.Fill) {
	publishFillEvent(ctx, s.orderBookStore, userId, fill)
}

func publishOrderEvent(ctx context.Context, store store.OrderBookStore, order *models.Order) {
	key, value, err := createOrderEvent(ctx, order)
	if err != nil {
		return
	}

	if err := store.PublishEvent(ctx, key, value); err != nil {
		logctx.Error(ctx, "failed to publish order event", logger.String("event", key), logger.Error(err))
	}
}

func (e *EvmClient) publishOrderEvent(ctx context.Context, order *models.Order) {
	publishOrderEvent(ctx, e.orderBookStore, order)
}
func (s *Service) publishOrderEvent(ctx context.Context, order *models.Order) {
	publishOrderEvent(ctx, s.orderBookStore, order)
}

func createOrderEvent(ctx context.Context, order *models.Order) (key string, value []byte, err error) {
	//value, err = order.ToJson()

	value, err = json.Marshal(struct {
		Event string `json:"event"`
		models.Order
	}{
		Event: "order-changed",
		Order: *order,
	})
	if err != nil {
		logctx.Error(ctx, "failed to marshal order to json", logger.Error(err))
	}

	key = models.CreateUserOrdersEventKey(order.UserId)

	return key, value, err
}
