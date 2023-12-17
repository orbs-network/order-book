package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (e *EvmClient) CheckPendingTxs(ctx context.Context) error {
	pendingSwaps, err := e.orderBookStore.GetPendingSwaps(ctx)
	if err != nil {
		logctx.Error(ctx, "Failed to get pending swaps", logger.Error(err))
		return err
	}
	fmt.Printf("pendingSwaps BEFORE: %#v\n", pendingSwaps)

	var wg sync.WaitGroup
	var mu sync.Mutex

	ptxs := make([]models.Pending, 0)

	for i := 0; i < len(pendingSwaps); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			p := pendingSwaps[index]

			tx, err := e.blockchainStore.GetTx(ctx, p.TxHash)
			if err != nil {
				if err == models.ErrNotFound {
					fmt.Println("Transaction not found")
					logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				} else {
					logctx.Error(ctx, "Failed to get transaction", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				}
				return
			}

			if tx == nil {
				fmt.Println("Transaction not found")
				logctx.Error(ctx, "Transaction not found but should be valid", logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
				return
			}

			switch tx.Status {
			case models.TX_SUCCESS:
				// Process successful transaction
				fmt.Println("Transaction successful  ----->")
				_, err = e.processSuccessfulTransaction(ctx, p, &mu)
				if err != nil {
					logctx.Error(ctx, "Failed to process successful transaction", logger.Error(err), logger.String("txHash", p.TxHash), logger.String("swapId", p.SwapId.String()))
					return
				}

			case models.TX_FAILURE:
				fmt.Printf("Transaction %s failed\n", p.TxHash)
				break
			case models.TX_PENDING:
				fmt.Printf("Transaction %s is still pending\n", p.TxHash)
				break
			default:
				fmt.Printf("Transaction %s has unknown status\n", p.TxHash)
				break
			}

		}(i)
	}

	wg.Wait() // Wait for all goroutines to complete

	mu.Lock()
	err = e.orderBookStore.StorePendingSwaps(ctx, ptxs)
	mu.Unlock()

	if err != nil {
		logctx.Error(ctx, "Failed to store updated pending swaps", logger.Error(err))
		return err
	}

	fmt.Println("Updated pending swaps list stored in Redis")
	return nil
}

func (e *EvmClient) processSuccessfulTransaction(ctx context.Context, p models.Pending, mu *sync.Mutex) ([]models.Order, error) {
	mu.Lock()
	defer mu.Unlock()

	orderFrags, err := e.orderBookStore.GetSwap(ctx, p.SwapId)
	if err != nil {
		logctx.Error(ctx, "Failed to get swap", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get swap: %w", err)
	}

	var orderIds []uuid.UUID
	for _, frag := range orderFrags {
		orderIds = append(orderIds, frag.OrderId)
	}

	orders, err := e.orderBookStore.FindOrdersByIds(ctx, orderIds)
	if err != nil {
		logctx.Error(ctx, "Failed to get orders", logger.Error(err), logger.String("swapId", p.SwapId.String()))
		return []models.Order{}, fmt.Errorf("failed to get orders: %w", err)
	}

	for _, order := range orders {
		isFilled, err := order.MarkSwapSuccess()
		if err != nil {
			logctx.Error(ctx, "Failed to mark order as filled", logger.Error(err), logger.String("orderId", order.Id.String()))
			return []models.Order{}, fmt.Errorf("failed to mark order as filled: %w", err)
		}

		if isFilled {
			logctx.Info(ctx, "Order is filled", logger.String("orderId", order.Id.String()))
			e.orderBookStore.StoreFilledOrders(ctx, []models.Order{order})
		} else {
			logctx.Info(ctx, "Order is partially filled", logger.String("orderId", order.Id.String()))
			e.orderBookStore.StoreOpenOrder(ctx, order)
		}
	}

	return orders, nil
}
