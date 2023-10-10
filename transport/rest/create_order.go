package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	Price  string `json:"price"`
	Size   string `json:"size"`
	Symbol string `json:"symbol"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	decPrice, err := decimal.NewFromString(order.Price)
	if err != nil {
		http.Error(w, "'price' is not a valid number format", http.StatusBadRequest)
		return
	}

	decSize, err := decimal.NewFromString(order.Size)
	if err != nil {
		http.Error(w, "'size' is not a valid number format", http.StatusBadRequest)
		return
	}

	symbol, err := models.StrToSymbol(order.Symbol)
	if err != nil {
		http.Error(w, "'symbol' is not a valid", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.svc.CreateOrder(r.Context(), decPrice, symbol, decSize)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedOrder)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
