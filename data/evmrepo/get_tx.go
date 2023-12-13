package evmrepo

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// GetTx returns the status of a transaction.
//
// If the transaction is pending, it returns TX_PENDING.
//
// If the transaction is successful, it returns TX_SUCCESS.
//
// If the transaction failed, it returns TX_FAILURE.
//
// If the transaction is not found, it returns a nil value and ErrNotFound.
func (e *evmRepository) GetTx(ctx context.Context, id string) (*models.Tx, error) {
	txHash := common.HexToHash(id)

	tx, pending, err := e.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		logctx.Error(ctx, "Error fetching transaction: %v", logger.Error(err), logger.String("txHash", txHash.String()))
		return nil, fmt.Errorf("error fetching transaction: %v", err)
	}

	if tx == nil {
		logctx.Error(ctx, "Transaction %q not found", logger.String("txHash", txHash.String()))
		return nil, models.ErrNotFound
	}

	if pending {
		logctx.Info(ctx, "Transaction %q is pending", logger.String("txHash", txHash.String()))
		return &models.Tx{
			Status: models.TX_PENDING,
			TxHash: tx.Hash().Hex(),
		}, nil
	}

	receipt, err := e.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		logctx.Error(ctx, "Error fetching transaction receipt: %v", logger.Error(err), logger.String("txHash", txHash.String()))
		return nil, fmt.Errorf("error fetching transaction receipt: %v", err)
	}

	if receipt.Status == 1 {
		logctx.Info(ctx, "Transaction %q succeeded", logger.String("txHash", txHash.String()))
		return &models.Tx{
			Status: models.TX_SUCCESS,
			TxHash: receipt.TxHash.Hex(),
		}, nil
	} else {
		logctx.Info(ctx, "Transaction %q failed", logger.String("txHash", txHash.String()))
		return &models.Tx{
			Status: models.TX_FAILURE,
			TxHash: receipt.TxHash.Hex(),
		}, nil
	}

}
