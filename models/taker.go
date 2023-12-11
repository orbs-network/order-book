package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BeginSwapRes struct {
	//Size       decimal.Decimal
	//OrderFrags []OrderFrag
	OutAmount decimal.Decimal
	Orders    []Order
	Fragments []OrderFrag
	SwapId    uuid.UUID
}

// type ConfirmAuctionRes struct {
// 	Orders        []models.Order
// 	Fragments     []models.OrderFrag
// 	BookSignature []byte
// }
