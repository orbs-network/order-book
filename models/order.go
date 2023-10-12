package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	STATUS_OPEN     Status = "OPEN"
	STATUS_PENDING  Status = "PENDING"
	STATUS_FILLED   Status = "FILLED"
	STATUS_CANCELED Status = "CANCELED"
)

type Order struct {
	Id        uuid.UUID       `json:"id"`
	UserId    uuid.UUID       `json:"userId"`
	Price     decimal.Decimal `json:"price"`
	Symbol    Symbol          `json:"symbol"`
	Size      decimal.Decimal `json:"size"`
	Signature *string         `json:"signature"` // EIP 712
	Status    Status          `json:"status"`    // when order is pending, it should not be updateable
}
