package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	Id            uuid.UUID       `json:"-"`
	UserId        uuid.UUID       `json:"userId"`
	Price         decimal.Decimal `json:"price"`
	Symbol        Symbol          `json:"symbol"`
	Size          decimal.Decimal `json:"size"`
	Signature     string          `json:"-" ` // EIP 712
	Status        Status          `json:"-"`  // when order is pending, it should not be updateable
	Side          Side            `json:"side"`
	Timestamp     time.Time       `json:"timestamp"`
	ClientOrderID uuid.UUID       `json:"clientOrderId"`
}

func (o *Order) OrderToMap() map[string]string {
	return map[string]string{
		"id":            o.Id.String(),
		"userId":        o.UserId.String(),
		"price":         o.Price.String(),
		"symbol":        o.Symbol.String(),
		"size":          o.Size.String(),
		"signature":     o.Signature,
		"status":        o.Status.String(),
		"side":          o.Side.String(),
		"timestamp":     o.Timestamp.String(),
		"clientOrderId": o.ClientOrderID.String(),
	}
}

func (o *Order) MapToOrder(data map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	idStr, exists := data["id"]
	if !exists {
		return fmt.Errorf("no id provided")
	}

	userIdStr, exists := data["userId"]
	if !exists {
		return fmt.Errorf("no userId provided")
	}

	priceStr, exists := data["price"]
	if !exists {
		return fmt.Errorf("no price provided")
	}

	symbolStr, exists := data["symbol"]
	if !exists {
		return fmt.Errorf("no symbol provided")
	}

	sizeStr, exists := data["size"]
	if !exists {
		return fmt.Errorf("no size provided")
	}

	signatureStr, exists := data["signature"]
	if !exists {
		return fmt.Errorf("no signature provided")
	}

	statusStr, exists := data["status"]
	if !exists {
		return fmt.Errorf("no status provided")
	}

	sideStr, exists := data["side"]
	if !exists {
		return fmt.Errorf("no side provided")
	}

	timestampStr, exists := data["timestamp"]
	if !exists {
		return fmt.Errorf("no timestamp provided")
	}

	clientOrderIdStr, exists := data["clientOrderId"]
	if !exists {
		return fmt.Errorf("no clientOrderId provided")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return fmt.Errorf("invalid userId: %v", err)
	}

	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		return fmt.Errorf("invalid price: %v", err)
	}

	size, err := decimal.NewFromString(sizeStr)
	if err != nil {
		return fmt.Errorf("invalid size: %v", err)
	}

	symbol, err := StrToSymbol(symbolStr)
	if err != nil {
		return err
	}

	side, err := StrToSide(sideStr)
	if err != nil {
		return err
	}

	status, err := StrToStatus(statusStr)
	if err != nil {
		return err
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %v", err)
	}

	clientOrderId, err := uuid.Parse(clientOrderIdStr)
	if err != nil {
		return fmt.Errorf("invalid clientOrderId: %v", err)
	}

	o.Id = id
	o.UserId = userId
	o.Price = price
	o.Symbol = symbol
	o.Size = size
	o.Signature = signatureStr
	o.Status = status
	o.Side = side
	o.Timestamp = timestamp
	o.ClientOrderID = clientOrderId

	return nil
}
