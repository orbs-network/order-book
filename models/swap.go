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

type QuoteRes struct {
	Size       decimal.Decimal
	OrderFrags []OrderFrag
}

type OrderFrag struct {
	OrderId uuid.UUID
	OutSize decimal.Decimal
	InSize  decimal.Decimal
}

func (f *OrderFrag) ToMap() map[string]string {
	return map[string]string{
		"inSize":  f.InSize.String(),
		"orderId": f.OrderId.String(),
		"outSize": f.OutSize.String(),
	}
}

func MarshalOrderFrags(orderFrags []OrderFrag) ([]byte, error) {
	swapMap := make([]map[string]string, len(orderFrags))
	for i, frag := range orderFrags {
		swapMap[i] = frag.ToMap()
	}

	return json.Marshal(swapMap)
}
