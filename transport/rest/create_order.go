package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	Price  string `json:"price"`
	Size   string `json:"size"`
	Symbol string `json:"symbol"`
}

// TODO: hardcoded userId for now
var userId = uuid.MustParse("d577273e-12de-4acc-a4f8-de7fb5b86e37")

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var args CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	decPrice, err := decimal.NewFromString(args.Price)
	if err != nil {
		http.Error(w, "'price' is not a valid number format", http.StatusBadRequest)
		return
	}

	decSize, err := decimal.NewFromString(args.Size)
	if err != nil {
		http.Error(w, "'size' is not a valid number format", http.StatusBadRequest)
		return
	}

	symbol, err := models.StrToSymbol(args.Symbol)
	if err != nil {
		http.Error(w, "'symbol' is not a valid", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to create order", logger.String("userId", userId.String()), logger.String("price", decPrice.String()), logger.String("size", decSize.String()))
	order, err := h.svc.AddOrder(r.Context(), userId, decPrice, symbol, decSize)

	if err == models.ErrOrderAlreadyExists {
		http.Error(w, "Order already exists", http.StatusConflict)
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
