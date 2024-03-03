package service

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func Key2Id(key string) string {
	splt := strings.Split(key, ":")

	if len(splt) < 2 {
		return ""
	}
	return splt[1]
}

func Key2UUID(key string) *uuid.UUID {
	id := Key2Id(key)
	if id == "" {
		return nil
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &uid
}

func (e *EvmClient) GetPendingSwaps(ctx context.Context) ([]models.Swap, error) {
	res := []models.Swap{}
	keys, err := e.orderBookStore.EnumSubKeysOf(ctx, "swapId")
	if err != nil {
		logctx.Error(ctx, "Failed to enum swapID keys", logger.Error(err))
		return res, err
	}
	for _, key := range keys {
		id := Key2UUID(key)
		if id != nil {
			swap, err := e.orderBookStore.GetSwap(ctx, *id)
			swap.Id = *id // confirm exists
			if err != nil {
				logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapKey", key))
			} else {
				// swap was started but not resolved
				if !swap.Started.IsZero() && swap.Resolved.IsZero() {
					// make sure id is there
					res = append(res, *swap)
				}
			}
		} else {
			logctx.Error(ctx, "failed to create id from key", logger.String("key", key))
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
		logctx.Info(ctx, "No pending swaps. Sleeping...")
		return nil
	}

	logctx.Info(ctx, "Found pending swaps to process", logger.Int("numPending", len(pendingSwaps)))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < len(pendingSwaps); i++ {
		logctx.Info(ctx, "Trying to process pending swap", logger.Int("index", i), logger.String("txHash", pendingSwaps[i].TxHash), logger.String("swapId", pendingSwaps[i].Id.String()))

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			p := pendingSwaps[index]

			tx, err := e.blockchainStore.GetTx(ctx, p.TxHash)
			if err != nil {
				if err == models.ErrNotFound {
					logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				} else {
					logctx.Error(ctx, "Failed to get transaction", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				}
				// abort failed tx
				return
			}

			if tx == nil {
				logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				return
			}

			switch tx.Status {
			case models.TX_SUCCESS:
				logctx.Info(ctx, "Transaction successful", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				_, err = e.ResolveSwap(ctx, p, true, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process successful transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
					return
				}
			case models.TX_FAILURE:
				logctx.Info(ctx, "Transaction failed", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				_, err = e.ResolveSwap(ctx, p, false, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process failed transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
					return
				}
			case models.TX_PENDING:
				logctx.Info(ctx, "Transaction still pending", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
			default:
				logctx.Error(ctx, "Unknown transaction status", logger.String("txHash", p.TxHash), logger.String("swapId", p.Id.String()))
				return
			}

		}(i)
	}

	wg.Wait()

	logctx.Info(ctx, "Finished checking pending swaps", logger.Int("numPending", len(pendingSwaps)))
	return nil
}
