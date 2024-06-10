package service

import (
	"context"
	"sync"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (e *EvmClient) GetPendingSwaps(ctx context.Context) ([]models.Swap, error) {
	res := []models.Swap{}
	openSwaps, err := e.orderBookStore.GetOpenSwaps(ctx)
	if err != nil {
		logctx.Error(ctx, "Failed to enum swap:open keys", logger.Error(err))
		return res, err
	}
	for _, swap := range openSwaps {

		// swap was started but not resolved
		if !swap.Started.IsZero() && swap.Resolved.IsZero() {
			// make sure id is there
			res = append(res, swap)
		}
	}
	return res, nil
}

// CheckPendingTxs checks all pending transactions and updates the order book accordingly
func (e *EvmClient) CheckPendingTxs(ctx context.Context) error {
	logctx.Debug(ctx, "Checking pending swaps...")
	pendingSwaps, err := e.GetPendingSwaps(ctx)
	if err != nil {
		logctx.Error(ctx, "Failed to get pending swaps", logger.Error(err))
		return err
	}

	if len(pendingSwaps) == 0 {
		logctx.Debug(ctx, "No pending swaps. Sleeping...")
		return nil
	}

	logctx.Debug(ctx, "Found pending swaps to process", logger.Int("numPending", len(pendingSwaps)))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < len(pendingSwaps); i++ {
		logctx.Debug(ctx, "Trying to process pending swap", logger.Int("index", i), logger.String("txHash", pendingSwaps[i].TxHash), logger.String("swapId", pendingSwaps[i].Id.String()))

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			p := pendingSwaps[index]

			tx, err := e.blockchainStore.GetTx(ctx, p.TxHash)
			if err != nil {
				if err == models.ErrNotFound {
					logctx.Warn(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				} else {
					logctx.Error(ctx, "Failed to get transaction", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				}
				// abort failed tx
				return
			}

			if tx == nil {
				logctx.Error(ctx, "Transaction is null", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				return
			}

			switch tx.Status {
			case models.TX_SUCCESS:
				logctx.Debug(ctx, "Transaction successful", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				p.Mined = *tx.Timestamp
				err = e.ResolveSwap(ctx, p, true, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process successful transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
					return
				}
			case models.TX_FAILURE:
				logctx.Debug(ctx, "Transaction failed", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				err = e.ResolveSwap(ctx, p, false, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process failed transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
					return
				}
			case models.TX_PENDING:
				logctx.Debug(ctx, "Transaction still pending", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
			default:
				logctx.Error(ctx, "Unknown transaction status", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				return
			}

		}(i)
	}

	wg.Wait()

	logctx.Debug(ctx, "Finished checking pending swaps", logger.Int("numPending", len(pendingSwaps)))
	return nil
}
