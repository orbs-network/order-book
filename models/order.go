package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	Id          uuid.UUID       `json:"orderId"`
	ClientOId   uuid.UUID       `json:"clientOrderId"`
	UserId      uuid.UUID       `json:"userId"`
	Price       decimal.Decimal `json:"price"`
	Symbol      Symbol          `json:"symbol"`
	Size        decimal.Decimal `json:"size"`
	SizePending decimal.Decimal `json:"-"`
	SizeFilled  decimal.Decimal `json:"-"`
	Signature   string          `json:"-" ` // EIP 712
	Side        Side            `json:"side"`
	Timestamp   time.Time       `json:"timestamp"`
}

func (o *Order) OrderToMap() map[string]string {
	return map[string]string{
		"id":          o.Id.String(),
		"clientOId":   o.ClientOId.String(),
		"userId":      o.UserId.String(),
		"price":       o.Price.String(),
		"symbol":      o.Symbol.String(),
		"size":        o.Size.String(),
		"sizePending": o.SizePending.String(),
		"sizeFilled":  o.SizeFilled.String(),
		"signature":   o.Signature,
		"side":        o.Side.String(),
		"timestamp":   o.Timestamp.Format(time.RFC3339),
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

	clientOIdStr, exists := data["clientOId"]
	if !exists {
		return fmt.Errorf("no clientOId provided")
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

	sizePendingStr, exists := data["sizePending"]
	if !exists {
		return fmt.Errorf("no sizePending provided")
	}

	sizeFilledStr, exists := data["sizeFilled"]
	if !exists {
		return fmt.Errorf("no sizeFilled provided")
	}

	signatureStr, exists := data["signature"]
	if !exists {
		return fmt.Errorf("no signature provided")
	}

	sideStr, exists := data["side"]
	if !exists {
		return fmt.Errorf("no side provided")
	}

	timestampStr, exists := data["timestamp"]
	if !exists {
		return fmt.Errorf("no timestamp provided")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("invalid id: %v", err)
	}

	clientOId, err := uuid.Parse(clientOIdStr)
	if err != nil {
		return fmt.Errorf("invalid clientOId: %v", err)
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

	sizePending, err := decimal.NewFromString(sizePendingStr)
	if err != nil {
		return fmt.Errorf("invalid sizePending: %v", err)
	}

	sizeFilled, err := decimal.NewFromString(sizeFilledStr)
	if err != nil {
		return fmt.Errorf("invalid sizeFilled: %v", err)
	}

	symbol, err := StrToSymbol(symbolStr)
	if err != nil {
		return err
	}

	side, err := StrToSide(sideStr)
	if err != nil {
		return err
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return fmt.Errorf("invalid timestamp: %v", err)
	}

	o.Id = id
	o.ClientOId = clientOId
	o.UserId = userId
	o.Price = price
	o.Symbol = symbol
	o.Size = size
	o.SizePending = sizePending
	o.SizeFilled = sizeFilled
	o.Signature = signatureStr
	o.Side = side
	o.Timestamp = timestamp

	return nil
}

// GetAvailableSize returns the size that is available to be filled
func (o *Order) GetAvailableSize() decimal.Decimal {
	used := o.SizePending.Add(o.SizeFilled)
	return o.Size.Sub(used)
}

// IsFilled returns true if the order has been filled
func (o *Order) IsFilled() bool {
	return o.SizePending.Equal(o.Size)
}

// IsPending returns true has a pending fill in progress
func (o *Order) IsPending() bool {
	return o.SizePending.GreaterThan(decimal.Zero)
}
