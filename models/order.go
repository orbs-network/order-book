package models

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/abi"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// EIP712 signature components
type Signature struct {
	Eip712Sig   string    `json:"sig"`
	AbiFragment abi.Order `json:"abiFragment"`
}

type Order struct {
	Id          uuid.UUID       `json:"orderId"`
	ClientOId   uuid.UUID       `json:"clientOrderId"`
	UserId      uuid.UUID       `json:"userId"`
	Price       decimal.Decimal `json:"price"`
	Symbol      Symbol          `json:"symbol"`
	Size        decimal.Decimal `json:"size"`
	SizePending decimal.Decimal `json:"pendingSize"`
	SizeFilled  decimal.Decimal `json:"filledSize"`
	Side        Side            `json:"side"`
	Timestamp   time.Time       `json:"timestamp"`
	Signature   Signature       `json:"-" `
	Cancelled   bool            `json:"cancelled"`
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
		"cancelled":   fmt.Sprintf("%t", o.Cancelled),
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

	cancelledStr, exists := data["cancelled"]
	if !exists {
		return fmt.Errorf("no cancelled provided")
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

	var abiFragment abi.Order

	if err := json.Unmarshal([]byte(abiFragmentJSON), &abiFragment); err != nil {
		return fmt.Errorf("invalid abiFragment: %v", err)
	}

	cancelled, err := strconv.ParseBool(cancelledStr)
	if err != nil {
		return fmt.Errorf("invalid cancelled value: %v", err)
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
	o.Cancelled = cancelled

	return nil
}

func (o *Order) ToJson() ([]byte, error) {
	return json.Marshal(o)
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

func (o *Order) IsUnfilled() bool {
	return o.SizeFilled.IsZero()
}

func (o *Order) IsPartialFilled() bool {
	return o.SizeFilled.GreaterThan(decimal.Zero) && o.SizeFilled.LessThan(o.Size)
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

func BigInt2Dcml(num *big.Int, dcmls int64) decimal.Decimal {
	// Convert the big.Int to decimal.Decimal
	decimalValue := decimal.NewFromBigInt(num, 0)

	// Normalize the decimal value by dividing by 10^18
	divisor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(dcmls))
	return decimalValue.Div(divisor)
}

// Status returns the status of the order
func (o *Order) OnchainPrice(inDec, outDec int) (decimal.Decimal, error) {
	// patch for tests with no onchain data
	if o.Signature.AbiFragment.Input.Amount == nil {
		return o.Price, nil
	}
	in := BigInt2Dcml(o.Signature.AbiFragment.Input.Amount, int64(inDec))
	fmt.Println("inAmount", in.String())
	out := BigInt2Dcml(o.Signature.AbiFragment.Outputs[0].Amount, int64(outDec))
	fmt.Println("outAmount", out.String())
	var result decimal.Decimal
	if o.Side == BUY {
		result = in.Div(out)
	} else {
		result = out.Div(in)
	}
	fmt.Println("onchain price: ", result.String())
	fmt.Println("origin  price: ", o.Price.String())

	return decimal.NewFromString(result.String())
}

func (o *Order) FragAtokenSize(frag OrderFrag) decimal.Decimal {
	if o.Side == BUY {
		// when maker buys, the A token is the Taker's inToken
		return frag.InSize
	} else {
		// when maker sells, the A token is the Taker's outToken
		return frag.OutSize
	}
}
func (o *Order) Fill(ctx context.Context, frag OrderFrag) (isFilled bool, err error) {
	fillSize := o.FragAtokenSize(frag)
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

// IMPORTANT!!!
// lock is always done in A TOKEN UNITs MATIC-USDC the lock is in MATIC and the size of the order is in MATIC
func (o *Order) Lock(ctx context.Context, frag OrderFrag) error {
	size := o.FragAtokenSize(frag)

	// cant lock cancelled liquidity
	if o.Cancelled {
		logctx.Error(ctx, "order is cancelled", logger.String("orderId", o.Id.String()))
		return ErrOrderCancelled
	}
	if o.GetAvailableSize().LessThan(size) {
		logctx.Error(ctx, "size to be locked greater than sizePending", logger.String("orderId", o.Id.String()), logger.String("pendingSize", o.SizePending.String()), logger.String("requestedLockSize", size.String()))
		return ErrUnexpectedSizeFilled
	}

	o.SizePending = o.SizePending.Add(size)
	return nil
}
func (o *Order) Unlock(ctx context.Context, frag OrderFrag) error {
	size := o.FragAtokenSize(frag)

	if o.SizePending.LessThan(size) {
		logctx.Error(ctx, "size to be unlocked is greater than sizePending", logger.String("orderId", o.Id.String()), logger.String("pendingSize", o.SizePending.String()), logger.String("requestedUnlockSize", size.String()))
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

// LogOrderDetails logs all the details of an order with the given message and log level
func LogOrderDetails(ctx context.Context, order *Order, msg string, level logctx.Level, err error, extraFields ...logger.Field) {
	fields := []logger.Field{
		logger.String("orderId", order.Id.String()),
		logger.String("clientOrderId", order.ClientOId.String()),
		logger.String("userId", order.UserId.String()),
		logger.String("symbol", string(order.Symbol)),
		logger.String("side", string(order.Side)),
		logger.String("price", order.Price.String()),
		logger.String("size", order.Size.String()),
		logger.String("sizePending", order.SizePending.String()),
		logger.String("sizeFilled", order.SizeFilled.String()),
		logger.String("timestamp", order.Timestamp.Format(time.RFC3339)),
		logger.String("cancelled", fmt.Sprintf("%t", order.Cancelled)),
	}

	if err != nil {
		fields = append(fields, logger.Error(err))
	}

	fields = append(fields, extraFields...)

	switch level {
	case logctx.DEBUG:
		logctx.Debug(ctx, msg, fields...)
	case logctx.INFO:
		logctx.Info(ctx, msg, fields...)
	case logctx.WARN:
		logctx.Warn(ctx, msg, fields...)
	case logctx.ERROR:
		logctx.Error(ctx, msg, fields...)
	}
}
