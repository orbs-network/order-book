package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	Price         string                 `json:"price"`
	Size          string                 `json:"size"`
	Symbol        string                 `json:"symbol"`
	Side          string                 `json:"side"`
	ClientOrderId string                 `json:"clientOrderId"`
	Eip712Sig     string                 `json:"eip712Sig"`
	Eip712MsgData map[string]interface{} `json:"eip712MsgData"`
}

type CreateOrderResponse struct {
	OrderId string `json:"orderId"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var args CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		logctx.Warn(ctx, "invalid JSON body", logger.Error(err))
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := handleValidateRequiredFields(hVRFArgs{
		price:         args.Price,
		size:          args.Size,
		symbol:        args.Symbol,
		side:          args.Side,
		clientOrderId: args.ClientOrderId,
		eip712Sig:     args.Eip712Sig,
		eip712MsgData: &args.Eip712MsgData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedFields, err := parseFields(w, pFInput{
		price:         args.Price,
		size:          args.Size,
		symbol:        args.Symbol,
		side:          args.Side,
		clientOrderId: args.ClientOrderId,
	})
	if err != nil {
		logctx.Warn(ctx, "failed to parse order fields", logger.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logctx.Info(ctx, "user trying to create order", logger.String("userId", user.Id.String()), logger.String("price", parsedFields.roundedDecPrice.String()), logger.String("size", parsedFields.decSize.String()), logger.String("clientOrderId", parsedFields.clientOrderId.String()))
	order, err := h.svc.CreateOrder(ctx, service.CreateOrderInput{
		UserId:        user.Id,
		Price:         parsedFields.roundedDecPrice,
		Symbol:        parsedFields.symbol,
		Size:          parsedFields.decSize,
		Side:          parsedFields.side,
		ClientOrderID: parsedFields.clientOrderId,
		Eip712Sig:     args.Eip712Sig,
		Eip712MsgData: args.Eip712MsgData,
	})

	if err == models.ErrSignatureVerificationError {
		logctx.Warn(ctx, "signature verification error", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, "Signature verification error", http.StatusBadRequest)
		return
	}

	if err == models.ErrSignatureVerificationFailed {
		logctx.Warn(ctx, "signature verification failed", logger.String("userId", user.Id.String()))
		http.Error(w, "Signature verification failed", http.StatusUnauthorized)
		return
	}

	if err == models.ErrClashingOrderId {
		http.Error(w, "Clashing order ID. Please retry", http.StatusConflict)
		return
	}

	if err == models.ErrClashingClientOrderId {
		http.Error(w, fmt.Sprintf("Order with clientOrderId %q already exists. You must first cancel this order", args.ClientOrderId), http.StatusConflict)
		return
	}

	if err != nil {
		logctx.Error(ctx, "failed to create order", logger.Error(err))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(CreateOrderResponse{
		OrderId: order.Id.String(),
	})

	if err != nil {
		logctx.Error(ctx, "failed to marshal created order", logger.Error(err))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err), logger.String("orderId", parsedFields.clientOrderId.String()))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
	}
}

type hVRFArgs struct {
	price         string
	size          string
	symbol        string
	side          string
	clientOrderId string
	eip712Sig     string
	eip712MsgData *map[string]interface{}
}

func handleValidateRequiredFields(args hVRFArgs) error {
	switch {
	case args.price == "":
		return fmt.Errorf("missing required field 'price'")

	case args.size == "":
		return fmt.Errorf("missing required field 'size'")

	case args.symbol == "":
		return fmt.Errorf("missing required field 'symbol'")

	case args.side == "":
		return fmt.Errorf("missing required field 'side'")

	case args.clientOrderId == "":
		return fmt.Errorf("missing required field 'clientOrderId'")

	case args.eip712Sig == "":
		return fmt.Errorf("missing required field 'eip712Sig'")

	case args.eip712MsgData == nil || *args.eip712MsgData == nil:
		return fmt.Errorf("missing required field 'eip712MsgData'")

	default:
		return nil
	}
}

type pfParsed struct {
	roundedDecPrice decimal.Decimal
	decSize         decimal.Decimal
	symbol          models.Symbol
	side            models.Side
	clientOrderId   uuid.UUID
}

type pFInput struct {
	price         string
	size          string
	symbol        string
	side          string
	clientOrderId string
}

func parseFields(w http.ResponseWriter, input pFInput) (*pfParsed, error) {
	decPrice, err := decimal.NewFromString(input.price)
	if err != nil {
		return nil, fmt.Errorf("'price' is not a valid number format")
	}
	if decPrice.IsNegative() {
		return nil, fmt.Errorf("'price' must be positive")
	}

	// TODO: Am I OK to always round?
	roundedDecPrice := decPrice.Round(2)

	decSize, err := decimal.NewFromString(input.size)
	if err != nil {
		return nil, fmt.Errorf("'size' is not a valid number format")
	}

	if decSize.IsNegative() {
		return nil, fmt.Errorf("'size' must be positive")
	}

	symbol, err := models.StrToSymbol(input.symbol)
	if err != nil {
		return nil, fmt.Errorf("'symbol' is not valid")
	}

	side, err := models.StrToSide(input.side)
	if err != nil {
		return nil, fmt.Errorf("'side' is not valid")
	}

	clientOrderId, err := uuid.Parse(input.clientOrderId)
	if err != nil {
		return nil, fmt.Errorf("'clientOrderId' is not valid")
	}

	return &pfParsed{
		roundedDecPrice: roundedDecPrice,
		decSize:         decSize,
		symbol:          symbol,
		side:            side,
		clientOrderId:   clientOrderId,
	}, nil
}
