package redisrepo

import (
	"context"

	"github.com/google/uuid"
)

func (r *redisRepository) RemoveOrder(ctx context.Context, orderId uuid.UUID) error {
	panic("not implemented")
}
