package service

import (
	"github.com/orbs-network/order-book/data/store"
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
		w.data[wallet] = decimal.Zero
	}
	w.data[wallet].Add(sum)
}

func (w *WalletVerifier) CheckAll(store *store.OrderBookStore) bool {
	return false
}
