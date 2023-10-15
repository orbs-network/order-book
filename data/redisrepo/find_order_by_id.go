package redisrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
)

func (r *redisRepository) FindOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error) {
	panic("not implemented")
}
