package models

import "errors"

type Symbol string

type strSet map[string]struct{}

// TODO: not final list of symbols

var (
	x          = struct{}{}
	symbolsMap = strSet{
		"MATIC-USDC":  x,
		"USDCE-USDT":  x,
		"ETH-BTC":     x,
		"ETH-USDCE":   x,
		"MATIC-ETH":   x,
		"MATIC-USDCE": x,
		"ETH-USDC":    x,
		"USDCE-USDC":  x,
		"ETH-USDT":    x,
		"DAI-USDCE":   x,
		"MATIC-USDT":  x,
		"BTC-USDCE":   x,
	}

	ErrInvalidSymbol = errors.New("invalid symbol")
)

func StrToSymbol(s string) (Symbol, error) {
	if _, exists := symbolsMap[s]; !exists {
		return "", ErrInvalidSymbol
	} else {
		return Symbol(s), nil
	}
}

func (s Symbol) String() string {
	return string(s)
}

func GetAllSymbols() []Symbol {
	symbols := make([]Symbol, 0, len(symbolsMap))
	for key := range symbolsMap {
		symbols = append(symbols, Symbol(key))
	}
	return symbols
}
