package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// EIP712 signature components
type Signature struct {
	Eip712Sig   string      `json:"sig"`
	AbiFragment AbiFragment `json:"abiFragment"`
}

type Order struct {
	Id        uuid.UUID       `json:"orderId"`
	ClientOId uuid.UUID       `json:"clientOrderId"`
	UserId    uuid.UUID       `json:"userId"`
	Price     decimal.Decimal `json:"price"`
	Symbol    Symbol          `json:"symbol"`
	Size      decimal.Decimal `json:"size"`
	// TODO: do we want to expose pending and filled sizes?
	SizePending decimal.Decimal `json:"pendingSize"`
	SizeFilled  decimal.Decimal `json:"filledSize"`
	Side        Side            `json:"side"`
	Timestamp   time.Time       `json:"timestamp"`
	Signature   Signature       `json:"-" `
}

func (o *Order) OrderToMap() map[string]string {
	// error can be ignored here because we know the data is valid
	abiFragmentBytes, _ := json.Marshal(o.Signature.AbiFragment)
	abiFragmentStr := string(abiFragmentBytes)

	return map[string]string{
		"id":          o.Id.String(),
		"clientOId":   o.ClientOId.String(),
		"userId":      o.UserId.String(),
		"price":       o.Price.String(),
		"symbol":      o.Symbol.String(),
		"size":        o.Size.String(),
		"sizePending": o.SizePending.String(),
		"sizeFilled":  o.SizeFilled.String(),
		"side":        o.Side.String(),
		"timestamp":   o.Timestamp.Format(time.RFC3339),
		"eip712Sig":   o.Signature.Eip712Sig,
		"abiFragment": abiFragmentStr,
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

	signatureStr, exists := data["eip712Sig"]
	if !exists {
		return fmt.Errorf("no signature provided")
	}

	abiFragmentJSON, exists := data["abiFragment"]
	if !exists {
		return fmt.Errorf("no abiFragmentJSON provided")
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

	var abiFragment AbiFragment

	if err := json.Unmarshal([]byte(abiFragmentJSON), &abiFragment); err != nil {
		return fmt.Errorf("invalid abiFragmen: %v", err)
	}

	o.Id = id
	o.ClientOId = clientOId
	o.UserId = userId
	o.Price = price
	o.Symbol = symbol
	o.Size = size
	o.SizePending = sizePending
	o.SizeFilled = sizeFilled
	o.Signature = Signature{
		Eip712Sig:   signatureStr,
		AbiFragment: abiFragment,
	}
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
	return o.SizeFilled.Equal(o.Size)
}

// IsPending returns true has a pending fill in progress
func (o *Order) IsPending() bool {
	return o.SizePending.GreaterThan(decimal.Zero)
}

// Status returns the status of the order
func (o *Order) Status() string {
	if o.IsFilled() {
		return "FILLED"
	}

	return "OPEN"
}

func (o *Order) Fill(ctx context.Context, fillSize decimal.Decimal) (isFilled bool, err error) {
	newSizeFilled := o.SizeFilled.Add(fillSize)
	if newSizeFilled.GreaterThan(o.Size) {
		logctx.Error(ctx, "total size is less than requested fill size", logger.String("orderId", o.Id.String()), logger.String("orderSize", o.Size.String()), logger.String("requestedFillSize", fillSize.String()))
		return false, ErrUnexpectedSizeFilled
	}

	if fillSize.GreaterThan(o.SizePending) {
		logctx.Error(ctx, "fillSize is greater than sizePending", logger.String("orderId", o.Id.String()), logger.String("pendingSize", o.SizePending.String()), logger.String("requestedFillSize", fillSize.String()))
		return false, ErrUnexpectedSizePending
	}

	o.SizeFilled = o.SizeFilled.Add(fillSize)
	o.SizePending = o.SizePending.Sub(fillSize)
	return o.IsFilled(), nil
}

func (o *Order) Kill(ctx context.Context, size decimal.Decimal) error {
	if o.SizePending.LessThan(size) {
		logctx.Error(ctx, "size to be rolled back is greater than sizePending", logger.String("orderId", o.Id.String()), logger.String("pendingSize", o.SizePending.String()), logger.String("requestedKillSize", size.String()))
		return ErrUnexpectedSizeFilled
	}

	o.SizePending = o.SizePending.Sub(size)
	return nil
}

// OrderIdsToStrings return a list of string order IDs from a list of orders
func OrderIdsToStrings(ctx context.Context, orders *[]Order) []string {
	if orders == nil {
		return []string{}
	}

	var orderIds []string
	for _, order := range *orders {
		orderIds = append(orderIds, order.Id.String())
	}
	return orderIds
}
