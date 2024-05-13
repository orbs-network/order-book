package service

import (
	"context"
	"math/big"
	"strings"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (e *EvmClient) UpdateMakerBalance(ctx context.Context, key string) {
	parts := strings.Split(key, ":")
	if len(parts) != 3 {
		logctx.Error(ctx, "CheckMakerBalance key invalid", logger.String("key", key))

	}
	// read db value
	val, err := e.orderBookStore.ReadStrKey(ctx, key)
	if err != nil {
		logctx.Error(ctx, "read balance key failed", logger.String("key", key))
		return
	}
	// current balance
	prvBlnc := new(big.Int)
	prvBlnc.SetString(val, 10)

	// update onchain value
	token := parts[1]
	maker := parts[2]

	blnc, err := e.blockchainStore.BalanceOf(ctx, token, maker)
	if err != nil {
		logctx.Error(ctx, "BalanceOf failed", logger.String("token", token), logger.String("maker", maker))
		return
	}

	// update balance if changed
	if blnc != prvBlnc {
		err := e.orderBookStore.WriteStrKey(ctx, key, blnc.String())
		if err != nil {
			logctx.Error(ctx, "WriteStrKey failed", logger.String("key", key), logger.String("blnc", blnc.String()))
		}
	}
}

func (e *EvmClient) UpdateMakerBalances(ctx context.Context) error {
	keys, err := e.orderBookStore.EnumSubKeysOf(ctx, "balance")
	if err != nil {
		logctx.Error(ctx, "EnumSubKeysOf balance failed")
		return err
	}
	for _, key := range keys {
		e.UpdateMakerBalance(ctx, key)
	}
	return nil
}
