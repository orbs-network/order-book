package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Input type for FindOrder method
type FindOrderInput struct {
	UserId uuid.UUID
	Price  decimal.Decimal
	Symbol Symbol
}
