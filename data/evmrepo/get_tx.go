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

	// First, try to get the transaction receipt
	receipt, err := e.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		// If error in fetching receipt, check if it's because the transaction is not found or pending
		tx, pending, errTx := e.client.TransactionByHash(ctx, txHash)

		if tx == nil && errTx.Error() == "not found" {
			logctx.Error(ctx, "Transaction %q not found", logger.String("txHash", txHash.String()))
			return nil, models.ErrNotFound
		}

		if errTx != nil {
			logctx.Error(ctx, "Error fetching transaction: %v", logger.Error(errTx), logger.String("txHash", txHash.String()))
			return nil, fmt.Errorf("error fetching transaction: %v", errTx)
		}

		if pending {
			logctx.Info(ctx, "Transaction %q is pending", logger.String("txHash", txHash.String()))
			return &models.Tx{
				Status: models.TX_PENDING,
				TxHash: tx.Hash().Hex(),
			}, nil
		}
	} else {
		// If receipt is found, check the status
		if receipt.Status == 1 {
			logctx.Info(ctx, "Transaction succeeded", logger.String("txHash", txHash.String()))
			return &models.Tx{
				Status: models.TX_SUCCESS,
				TxHash: receipt.TxHash.Hex(),
			}, nil
		} else {
			logctx.Info(ctx, "Transaction failed", logger.String("txHash", txHash.String()))
			return &models.Tx{
				Status: models.TX_FAILURE,
				TxHash: receipt.TxHash.Hex(),
			}, nil
		}
	}

	// Unexpected case
	logctx.Error(ctx, "Unexpected case for transaction %q", logger.String("txHash", txHash.String()))
	return nil, fmt.Errorf("unexpected case for transaction %q", txHash.String())
}
