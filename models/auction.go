package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AuctionStatus string

const (
	AUCTION_CONFIRMED AuctionStatus = "confirmed"
	AUCTION_MINED     AuctionStatus = "mined"
	AUCTION_REVERTED  AuctionStatus = "reverted"
)

func (a AuctionStatus) String() string {
	return string(a)
}

type AmountOut struct {
	Size       decimal.Decimal
	OrderFrags []OrderFrag
}

type OrderFrag struct {
	OrderId uuid.UUID
	Size    decimal.Decimal
}

func (f *OrderFrag) ToMap() map[string]string {
	return map[string]string{
		"orderId": f.OrderId.String(),
		"size":    f.Size.String(),
	}
}

func (f *OrderFrag) ToOrderFrag(data map[string]string) error {
	if len(data) == 0 {
		return nil
	}

	orderIdStr, exists := data["orderId"]
	if !exists {
		return nil
	}

	amountStr, exists := data["size"]
	if !exists {
		return nil
	}

	orderId, err := uuid.Parse(orderIdStr)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return err
	}

	f.OrderId = orderId
	f.Size = amount

	return nil
}

func MarshalOrderFrags(orderFrags []OrderFrag) ([]byte, error) {
	auctionMap := make([]map[string]string, len(orderFrags))
	for i, frag := range orderFrags {
		auctionMap[i] = frag.ToMap()
	}

	return json.Marshal(auctionMap)
}
