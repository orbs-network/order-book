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
	prvBlnc := new(big.Float)
	prvBlnc.SetString(val)

	// update onchain value
	token := parts[1]
	maker := parts[2]

	blnc, err := e.blockchainStore.BalanceOf(ctx, token, maker)
	if err != nil {
		logctx.Error(ctx, "BalanceOf failed", logger.String("token", token), logger.String("maker", maker))
		return
	}

	// normalize balance to size with decimals
	dcmls, err := e.blockchainStore.TokenDecimals(ctx, token, maker)
	if err != nil {
		logctx.Error(ctx, "TokenDecimals failed", logger.String("token", token))
		return
	}
	// Create a big.Int representing 10^exponent
	denom := new(big.Int).Exp(big.NewInt(10), big.NewInt(dcmls), nil)

	fBlnc := new(big.Float).Quo(
		new(big.Float).SetInt(blnc),
		new(big.Float).SetInt(denom),
	)

	// update balance if changed
	if fBlnc != prvBlnc {
		err := e.orderBookStore.WriteStrKey(ctx, key, fBlnc.String())
		if err != nil {
			logctx.Error(ctx, "WriteStrKey failed", logger.String("key", key), logger.String("blnc", fBlnc.String()))
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
