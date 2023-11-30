package models

import (
	"github.com/google/uuid"
)

type BeginSwapRes struct {
	//Size       decimal.Decimal
	//OrderFrags []OrderFrag
	Orders    []Order
	Fragments []OrderFrag
	SwapId    uuid.UUID
}

// type ConfirmAuctionRes struct {
// 	Orders        []models.Order
// 	Fragments     []models.OrderFrag
// 	BookSignature []byte
// }
