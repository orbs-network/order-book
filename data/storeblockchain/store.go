package storeblockchain

import (
	"context"
	"math/big"

	"github.com/orbs-network/order-book/models"
)

type BlockchainStore interface {
	GetTx(ctx context.Context, id string) (*models.Tx, error)
	BalanceOf(ctx context.Context, token, adrs string) (*big.Int, error)
	TokenDecimals(ctx context.Context, token, adrs string) (int64, error)
}
