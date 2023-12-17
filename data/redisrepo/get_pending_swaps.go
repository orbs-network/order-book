package redisrepo

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
)

// GetPendingSwaps returns all pending swaps that are waiting to be checked for completion
func (r *redisRepository) GetPendingSwaps(ctx context.Context) ([]models.Pending, error) {
	var pendingSwaps []models.Pending

	pendings, err := r.client.LRange(ctx, CreatePendingSwapTxsKey(), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending swaps: %s", err)
	}

	for _, pending := range pendings {
		var p models.Pending
		err := p.FromJson([]byte(pending))
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal pending: %s", err)
		}

		pendingSwaps = append(pendingSwaps, p)
	}

	return pendingSwaps, nil
}
