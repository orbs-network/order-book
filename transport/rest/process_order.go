package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	Price         string `json:"price"`
	Size          string `json:"size"`
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`
	ClientOrderId string `json:"clientOrderId"`
}

// TODO: hardcoded userId for now
var userId = uuid.MustParse("d577273e-12de-4acc-a4f8-de7fb5b86e37")

func (h *Handler) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	var args CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := handleValidateRequiredFields(args); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decPrice, err := decimal.NewFromString(args.Price)
	if err != nil {
		http.Error(w, "'price' is not a valid number format", http.StatusBadRequest)
		return
	}
	if decPrice.IsNegative() {
		http.Error(w, "'price' must be positive", http.StatusBadRequest)
		return
	}

	// TODO: Am I OK to always round?
	roundedDecPrice := decPrice.Round(2)

	decSize, err := decimal.NewFromString(args.Size)
	if err != nil {
		http.Error(w, "'size' is not a valid number format", http.StatusBadRequest)
		return
	}

	if !decSize.IsInteger() {
		http.Error(w, "'size' must be an integer", http.StatusBadRequest)
		return
	}

	if decSize.IsNegative() {
		http.Error(w, "'size' must be positive", http.StatusBadRequest)
		return
	}

	symbol, err := models.StrToSymbol(args.Symbol)
	if err != nil {
		http.Error(w, "'symbol' is not valid", http.StatusBadRequest)
		return
	}

	side, err := models.StrToSide(args.Side)
	if err != nil {
		http.Error(w, "'side' is not valid", http.StatusBadRequest)
		return
	}

	clientOrderId, err := uuid.Parse(args.ClientOrderId)
	if err != nil {
		http.Error(w, "'clientOrderId' is not valid", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to create order", logger.String("userId", userId.String()), logger.String("price", roundedDecPrice.String()), logger.String("size", decSize.String()), logger.String("clientOrderId", clientOrderId.String()))
	order, err := h.svc.ProcessOrder(r.Context(), service.ProcessOrderInput{
		UserId:        userId,
		Price:         roundedDecPrice,
		Symbol:        symbol,
		Size:          decSize,
		Side:          side,
		ClientOrderID: clientOrderId,
	})

	if err == models.ErrOrderAlreadyExists {
		http.Error(w, "Order already exists. You must first cancel existing order", http.StatusConflict)
		return
	}

	if err == service.ErrClashingOrderId {
		http.Error(w, "Clashing 'clientOrderId'. Retry with a different UUID", http.StatusConflict)
		return
	}

	if err != nil {
		logctx.Error(r.Context(), "failed to create order", logger.Error(err))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(order)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal created order", logger.Error(err))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func handleValidateRequiredFields(args CreateOrderRequest) error {
	switch {
	case args.Price == "":
		return fmt.Errorf("missing required field 'price'")

	case args.Size == "":
		return fmt.Errorf("missing required field 'size'")

	case args.Symbol == "":
		return fmt.Errorf("missing required field 'symbol'")

	case args.Side == "":
		return fmt.Errorf("missing required field 'side'")

	case args.ClientOrderId == "":
		return fmt.Errorf("missing required field 'clientOrderId'")

	default:
		return nil
	}
}
