package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// CheckPendingTxs checks all pending transactions and updates the order book accordingly
func (e *EvmClient) CheckPendingTxs(ctx context.Context) error {
	logctx.Info(ctx, "Checking pending transactions...")
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
				_, err = e.processCompletedTransaction(ctx, p, true, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process successful transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
					return
				}
			case models.TX_FAILURE:
				logctx.Info(ctx, "Transaction failed", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				_, err = e.processCompletedTransaction(ctx, p, false, &mu)
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
	err = e.orderBookStore.StorePendingSwaps(ctx, ptxs)
	mu.Unlock()

	if err != nil {
		logctx.Error(ctx, "Failed to store updated pending swaps", logger.Error(err))
		return err
	}

	logctx.Info(ctx, "Finished checking pending transactions", logger.Int("numPending", len(ptxs)))
	return nil
}

func (e *EvmClient) processCompletedTransaction(ctx context.Context, p models.SwapTx, isSuccessful bool, mu *sync.Mutex) ([]models.Order, error) {
	mu.Lock()
	defer mu.Unlock()

	orderFrags, err := e.orderBookStore.GetSwap(ctx, p.SwapId)
	if err != nil {
		logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get swap: %w", err)
	}

	var orderIds []uuid.UUID
	orderSizes := make(map[uuid.UUID]decimal.Decimal)

	for _, frag := range orderFrags {
		orderIds = append(orderIds, frag.OrderId)
		orderSizes[frag.OrderId] = frag.OutSize
	}

	orders, err := e.orderBookStore.FindOrdersByIds(ctx, orderIds, false)
	if err != nil {
		logctx.Error(ctx, "Failed to get orders", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get orders: %w", err)
	}

	var swapOrders []*models.Order

	for _, order := range orders {

		size, found := orderSizes[order.Id]
		if !found {
			logctx.Error(ctx, "Failed to get order frag size", logger.String("orderId", order.Id.String()))
			return []models.Order{}, fmt.Errorf("failed to get order frag size")
		}

		if isSuccessful {
			if _, err := order.Fill(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to mark order as filled", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}
		} else {
			if err := order.Unlock(ctx, size); err != nil {
				logctx.Error(ctx, "Failed to Release order locked liq", logger.Error(err), logger.String("orderId", order.Id.String()))
				continue
			}
		}
		swapOrders = append(swapOrders, &order)
	}

	err = e.orderBookStore.ProcessCompletedSwapOrders(ctx, swapOrders, p.SwapId, isSuccessful)
	if err != nil {
		logctx.Error(ctx, "Failed to process completed swap orders", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return []models.Order{}, fmt.Errorf("failed to process completed swap orders: %w", err)
	}

	return orders, nil
}
