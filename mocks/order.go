package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

var OrderId = uuid.MustParse("00000000-0000-0000-0000-000000000008")
var ClientOId = uuid.MustParse("00000000-0000-0000-0000-000000000009")
var UserId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var Price, _ = decimal.NewFromString("0.85450000")
var Symbol, _ = models.StrToSymbol("MATIC-USDC")
var Size, _ = decimal.NewFromString("126")
var Side = models.BUY
var Timestamp = time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)

var Order = models.Order{
	Id:        OrderId,
	ClientOId: ClientOId,
	UserId:    UserId,
	Price:     Price,
	Symbol:    Symbol,
	Size:      Size,
	Side:      Side,
	Timestamp: Timestamp,
}
