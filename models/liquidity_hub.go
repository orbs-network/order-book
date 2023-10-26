package models

import (
	"encoding/json"

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

func (f *FilledOrder) FilledOrderToMap() map[string]string {
	return map[string]string{
		"orderId": f.OrderId.String(),
		"amount":  f.Amount.String(),
	}
}

func (f *FilledOrder) MapToFilledOrder(data map[string]string) error {
	if len(data) == 0 {
		return nil
	}

	orderIdStr, exists := data["orderId"]
	if !exists {
		return nil
	}

	amountStr, exists := data["amount"]
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
	f.Amount = amount

	return nil
}

func MarshalFilledOrders(fillOrders []FilledOrder) ([]byte, error) {
	auctionMap := make([]map[string]string, len(fillOrders))
	for i, fillOrder := range fillOrders {
		auctionMap[i] = fillOrder.FilledOrderToMap()
	}

	return json.Marshal(auctionMap)
}
