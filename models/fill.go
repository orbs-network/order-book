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
	Mined     time.Time       `json:"mined"`
	Resolved  time.Time       `json:"resolved"`
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	OrderSize decimal.Decimal `json:"orderSize"`
}

func NewFill(symbol Symbol, swap Swap, frag OrderFrag, order *Order) *Fill {
	// create fill res
	fill := Fill{
		OrderId:  frag.OrderId,
		SwapId:   swap.Id,
		Symbol:   symbol,
		Mined:    swap.Mined,
		Resolved: swap.Resolved,
	}

	// get order
	if order != nil {
		fill.Size = order.FragAtokenSize(frag)

		// enrich fill data
		fill.Side = order.Side
		fill.ClientOId = order.ClientOId
		fill.Price = order.Price
		fill.OrderSize = order.Size
	}
	return &fill
}
