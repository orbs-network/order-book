package models

import "errors"

type Symbol string

// TODO: not final list of symbols
var (
	symbolsMap = map[string]Symbol{
		"BTC-ETH":  "BTC-ETH",
		"ETH-BTC":  "ETH-BTC",
		"USDC-ETH": "USDC-ETH",
		"ETH-USDC": "ETH-USDC",
		"ETH-USD":  "ETH-USD",
		"USD-ETH":  "USD-ETH",
		"TRX-BTT":  "TRX-BTT",
		"BTT-TRX":  "BTT-TRX",
		"BTC-USDC": "BTC-USDC",
		"USDC-BTC": "USDC-BTC",
		"BTC-USD":  "BTC-USD",
		"USD-BTC":  "USD-BTC",
		"TRX-ETH":  "TRX-ETH",
		"ETH-TRX":  "ETH-TRX",
		"BTT-ETH":  "BTT-ETH",
		"ETH-BTT":  "ETH-BTT",
		"TRX-USD":  "TRX-USD",
		"USD-TRX":  "USD-TRX",
		"BTT-USD":  "BTT-USD",
		"USD-BTT":  "USD-BTT",
		"TRX-USDC": "TRX-USDC",
		"USDC-TRX": "USDC-TRX",
		"BTT-USDC": "BTT-USDC",
		"USDC-BTT": "USDC-BTT",
		"TRX-BTC":  "TRX-BTC",
		"BTC-TRX":  "BTC-TRX",
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
