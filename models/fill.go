package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Individual Fill
// get SwapFills for user
type Fill struct {
	OrderId   uuid.UUID       `json:"orderId"`
	ClientOId uuid.UUID       `json:"clientOrderId"`
	SwapId    uuid.UUID       `json:"swapId"`
	Side      Side            `json:"side"`
	Symbol    Symbol          `json:"symbol"`
	Timestamp time.Time       `json:"timestamp"`
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	OrderSize decimal.Decimal `json:"orderSize"`

	Cancelled bool `json:"cancelled"`
}
