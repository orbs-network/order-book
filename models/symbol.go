package models

import "errors"

type Symbol string

// TODO: not final list of symbols
var (
	symbolsMap = map[string]Symbol{
		"MATIC-USDC":  "MATIC-USDC",
		"USDCE-USDT":  "USDCE-USDT",
		"ETH-BTC":     "ETH-BTC",
		"ETH-USDCE":   "ETH-USDCE",
		"MATIC-ETH":   "MATIC-ETH",
		"MATIC-USDCE": "MATIC-USDCE",
		"ETH-USDC":    "ETH-USDC",
		"USDCE-USDC":  "USDCE-USDC",
		"ETH-USDT":    "ETH-USDT",
		"DAI-USDCE":   "DAI-USDCE",
		"MATIC-USDT":  "MATIC-USDT",
		"BTC-USDCE":   "BTC-USDCE",
	}

	ErrInvalidSymbol = errors.New("invalid symbol")
)

func StrToSymbol(s string) (Symbol, error) {
	if symbol, ok := symbolsMap[s]; ok {
		return symbol, nil
	}
	return "", ErrInvalidSymbol
}

func (s Symbol) String() string {
	return string(s)
}

func GetAllSymbols() []Symbol {
	symbols := make([]Symbol, 0, len(symbolsMap))
	for _, symbol := range symbolsMap {
		symbols = append(symbols, symbol)
	}
	return symbols
}
