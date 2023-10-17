package models

import "github.com/shopspring/decimal"

type MarketDepth struct {
	Asks   [][]decimal.Decimal `json:"asks"`
	Bids   [][]decimal.Decimal `json:"bids"`
	Symbol string              `json:"symbol"`
	Time   int64               `json:"time"`
}
