package models

import "errors"

type Symbol string

const (
	BTC_ETH  Symbol = "BTC-ETH"
	USDC_ETH Symbol = "USDC-ETH"
	ETH_USD  Symbol = "ETH-USD"
)

var ErrInvalidSymbol = errors.New("invalid symbol")

func StrToSymbol(s string) (Symbol, error) {
	switch s {
	case "BTC-ETH":
		return BTC_ETH, nil
	case "USDC-ETH":
		return USDC_ETH, nil
	case "ETH-USD":
		return ETH_USD, nil
	// TODO: add more symbols
	default:
		return "", ErrInvalidSymbol
	}
}

func (s Symbol) String() string {
	return string(s)
}
