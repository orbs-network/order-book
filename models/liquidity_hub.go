package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AmountOut struct {
	AmountOut  decimal.Decimal
	FillOrders []FilledOrder
}

type FilledOrder struct {
	OrderId uuid.UUID
	Amount  decimal.Decimal
}
