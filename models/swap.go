package models

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SwapStatus string

const (
	// taker
	SWAP_STARTED  SwapStatus = "swap_started"
	SWAP_ABORDTED SwapStatus = "swap_aborted"
)

func (a SwapStatus) String() string {
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
	swapMap := make([]map[string]string, len(orderFrags))
	for i, frag := range orderFrags {
		swapMap[i] = frag.ToMap()
	}

	return json.Marshal(swapMap)
}
