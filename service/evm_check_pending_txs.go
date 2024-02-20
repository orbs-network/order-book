package service

import (
	"context"
	"sync"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// CheckPendingTxs checks all pending transactions and updates the order book accordingly
func (e *EvmClient) CheckPendingTxs(ctx context.Context) error {
	logctx.Debug(ctx, "Checking pending transactions...")
	pendingSwaps, err := e.orderBookStore.GetPendingSwaps(ctx)
	if err != nil {
		logctx.Error(ctx, "Failed to get pending swaps", logger.Error(err))
		return err
	}

	if len(pendingSwaps) == 0 {
		logctx.Info(ctx, "No pending transactions. Sleeping...")
		return nil
	}

	logctx.Info(ctx, "Found pending transactions to process", logger.Int("numPending", len(pendingSwaps)))

	var wg sync.WaitGroup
	var mu sync.Mutex

	ptxs := make([]models.SwapTx, 0)

	for i := 0; i < len(pendingSwaps); i++ {

		logctx.Debug(ctx, "Trying to process pending transaction", logger.Int("index", i), logger.String("txHash", pendingSwaps[i].TxHash), logger.String("swapId", pendingSwaps[i].SwapId.String()))

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			p := pendingSwaps[index]

			tx, err := e.blockchainStore.GetTx(ctx, p.TxHash)
			if err != nil {
				if err == models.ErrNotFound {
					logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				} else {
					logctx.Error(ctx, "Failed to get transaction", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				}
				return
			}

			if tx == nil {
				logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				return
			}

			switch tx.Status {
			case models.TX_SUCCESS:
				logctx.Info(ctx, "Transaction successful", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				_, err = e.ProcessCompletedTransaction(ctx, p, true, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process successful transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
					return
				}
			case models.TX_FAILURE:
				logctx.Info(ctx, "Transaction failed", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				_, err = e.ProcessCompletedTransaction(ctx, p, false, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process failed transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
					return
				}
			case models.TX_PENDING:
				logctx.Info(ctx, "Transaction still pending", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				ptxs = append(ptxs, p)
			default:
				logctx.Error(ctx, "Unknown transaction status", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				return
			}

		}(i)
	}

	wg.Wait()

	mu.Lock()
	// Store swaps that are still pending
	err = e.orderBookStore.StorePendingSwaps(ctx, ptxs)
	mu.Unlock()

	if err != nil {
		logctx.Error(ctx, "Failed to store updated pending swaps", logger.Error(err))
		return err
	}

	logctx.Info(ctx, "Finished checking pending transactions", logger.Int("numPending", len(ptxs)))
	return nil
}
