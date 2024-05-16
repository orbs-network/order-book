package service

import (
	"context"

	"github.com/orbs-network/order-book/data/store"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type wallet2Sum map[string]decimal.Decimal
type WalletVerifier struct {
	tokenAdrs string
	data      wallet2Sum
}

func NewWalletVerifier(tokenAdrs string) *WalletVerifier {
	return &WalletVerifier{
		tokenAdrs: tokenAdrs,
		data:      make(wallet2Sum),
	}
}

func (w *WalletVerifier) Add(wallet string, sum decimal.Decimal) {
	_, exists := w.data[wallet]
	// add entry
	if !exists {
		w.data[wallet] = decimal.NewFromInt(0)
	}
	w.data[wallet] = w.data[wallet].Add(sum)
}

func (w *WalletVerifier) CheckOne(ctx context.Context, st store.OrderBookStore, wallet string, sum decimal.Decimal) bool {
	blnc, err := st.GetMakerTokenBalance(ctx, w.tokenAdrs, wallet)
	if err != nil {
		logctx.Error(ctx, "ReadStrKey failed", logger.String("token", w.tokenAdrs), logger.String("wallet", wallet), logger.Error(err))
		return false
	}
	logctx.Info(ctx, "QuoteVsBalance", logger.String("token", w.tokenAdrs), logger.String("wallet", wallet), logger.String("quoteSize", sum.String()), logger.String("balance", blnc.String()))
	return sum.LessThanOrEqual(blnc)
}
func (w *WalletVerifier) CheckAll(ctx context.Context, st store.OrderBookStore) bool {
	for wallet, sum := range w.data {
		if !w.CheckOne(ctx, st, wallet, sum) {
			return false
		}
	}
	return true
}
