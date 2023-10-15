package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

func (r *redisRepository) FindOrder(ctx context.Context, input models.FindOrderInput) (*models.Order, error) {
	panic("not implemented")
}
