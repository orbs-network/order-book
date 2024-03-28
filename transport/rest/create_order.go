package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/restutils"
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
	Eip712Msg     map[string]interface{} `json:"eip712Msg"`
}

type CreateOrderResponse struct {
	OrderId string `json:"orderId"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	var args CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		logctx.Warn(ctx, "invalid JSON body", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if err := handleValidateRequiredFields(hVRFArgs{
		price:         args.Price,
		size:          args.Size,
		symbol:        args.Symbol,
		side:          args.Side,
		clientOrderId: args.ClientOrderId,
		eip712Sig:     args.Eip712Sig,
		eip712Msg:     &args.Eip712Msg,
	}); err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error())
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
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	abiFragment, err := restutils.ConvertToAbiFragment(args.Eip712Msg)
	if err != nil {
		logctx.Warn(ctx, "failed to convert eip712Msg to abi fragment", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, fmt.Errorf("failed to parse eip712Msg: %w", err).Error())
		return
	}

	logctx.Debug(ctx, "user trying to create order", logger.String("userId", user.Id.String()), logger.String("price", parsedFields.roundedDecPrice.String()), logger.String("size", parsedFields.decSize.String()), logger.String("clientOrderId", parsedFields.clientOrderId.String()))
	order, err := h.svc.CreateOrder(ctx, service.CreateOrderInput{
		UserId:        user.Id,
		Price:         parsedFields.roundedDecPrice,
		Symbol:        parsedFields.symbol,
		Size:          parsedFields.decSize,
		Side:          parsedFields.side,
		ClientOrderID: parsedFields.clientOrderId,
		Eip712Sig:     args.Eip712Sig,
		AbiFragment:   abiFragment,
	})

	if err == models.ErrClashingOrderId {
		logctx.Warn(ctx, "clashing order ID", logger.String("userId", user.Id.String()), logger.String("orderId", parsedFields.clientOrderId.String()))
		restutils.WriteJSONError(ctx, w, http.StatusConflict, "Clashing order ID. Please retry")
		return
	}

	if err == models.ErrClashingClientOrderId {
		logctx.Warn(ctx, "clashing client order ID", logger.String("userId", user.Id.String()), logger.String("clientOrderId", parsedFields.clientOrderId.String()))
		restutils.WriteJSONError(ctx, w, http.StatusConflict, fmt.Sprintf("Order with clientOrderId %s already exists", args.ClientOrderId))
		return
	}

	if err != nil {
		logctx.Error(ctx, "failed to create order", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error creating order. Try again later")
		return
	}

	resp, err := json.Marshal(CreateOrderResponse{
		OrderId: order.Id.String(),
	})

	if err != nil {
		logctx.Error(ctx, "failed to marshal created order", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error creating order. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err), logger.String("orderId", parsedFields.clientOrderId.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error creating order. Try again later")
	}

	logctx.Debug(ctx, "order created", logger.String("userId", user.Id.String()), logger.String("orderId", order.Id.String()), logger.String("price", parsedFields.roundedDecPrice.String()), logger.String("size", parsedFields.decSize.String()), logger.String("side", parsedFields.side.String()))
}

type hVRFArgs struct {
	price         string
	size          string
	symbol        string
	side          string
	clientOrderId string
	eip712Sig     string
	eip712Msg     *map[string]interface{}
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

	case args.eip712Msg == nil || *args.eip712Msg == nil:
		return fmt.Errorf("missing required field 'eip712Msg'")

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

func parseFields(_ http.ResponseWriter, input pFInput) (*pfParsed, error) {
	decPrice, err := decimal.NewFromString(input.price)
	if err != nil {
		return nil, fmt.Errorf("'price' is not a valid number format")
	}

	if decPrice.IsZero() || decPrice.IsNegative() {
		return nil, fmt.Errorf("'price' must be positive")
	}

	// Ensure price is 8 decimal places - aligns with Binance's precision
	if decPrice.Exponent() < -8 {
		return nil, fmt.Errorf("'price' must not exceed 8 decimal places")
	}
	roundedDecPrice := decPrice.Round(8)

	decSize, err := decimal.NewFromString(input.size)
	if err != nil {
		return nil, fmt.Errorf("'size' is not a valid number format")
	}

	if decSize.IsZero() || decSize.IsNegative() {
		return nil, fmt.Errorf("'size' must be positive")
	}

	if decSize.Exponent() < -4 {
		return nil, fmt.Errorf("'size' must not exceed 4 decimal places")
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
