package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (e *EvmClient) UpdateMakerBalance(ctx context.Context) error {
	keys, err := e.orderBookStore.EnumSubKeysOf(ctx, "balance")
	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) != 3 {
			logctx.Error(ctx, "CheckMakerBalance key invalid", logger.String("key", key))
			return err

		}
		// read db value

		// update onchain value
		token := parts[1]
		maker := parts[2]

		blnc, err := e.blockchainStore.BalanceOf(ctx, token, maker)
		if err != nil {
			return err
		}

		// ipdate balance if changed
		fmt.Println(blnc)

	}
	return nil
}
