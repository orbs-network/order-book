package memoryrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (r *inMemoryRepository) FindOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	orderIdStr := orderId.String()

	element, exists := r.orderLocations[orderIdStr]
	if !exists {
		logctx.Info(ctx, "Order not found", logger.String("orderId", orderIdStr))
		return nil, nil
	}

	order, ok := element.Value.(models.Order)
	if !ok {
		logctx.Error(ctx, "Failed to cast order", logger.String("orderId", orderIdStr))
		return nil, fmt.Errorf("failed to cast order for orderId %q", orderIdStr)
	}

	return &order, nil
}
