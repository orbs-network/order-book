package models

import "errors"

type Symbol string

// TODO: not final list of symbols
var (
	symbolsMap = map[string]Symbol{
		"MATIC-USDC": "MATIC-USDC",
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
