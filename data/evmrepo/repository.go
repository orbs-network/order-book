package evmrepo

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type evmRepository struct {
	client Blockchain
}

func NewEvmRepository(client Blockchain) (*evmRepository, error) {
	return &evmRepository{
		client: client,
	}, nil
}

type Blockchain interface {
	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// TransactionByHash returns the transaction by transaction hash.
	TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
	// BlockByNumber returns the block by block number.
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}
