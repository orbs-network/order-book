package storeblockchain

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

type BlockchainStore interface {
	GetTx(ctx context.Context, id string) (*models.Tx, error)
}
