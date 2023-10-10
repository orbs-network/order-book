package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	Id        uuid.UUID
	Price     decimal.Decimal
	Symbol    Symbol
	Size      decimal.Decimal
	Signature *string // EIP 712
	Pending   bool    // when order is pending, it should not be updateable
}
