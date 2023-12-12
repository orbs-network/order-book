package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BeginSwapRes struct {
	OutAmount decimal.Decimal
	Orders    []Order
	Fragments []OrderFrag
	SwapId    uuid.UUID
}
