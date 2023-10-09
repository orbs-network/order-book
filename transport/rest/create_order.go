package rest

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.svc.CreateOrder(r.Context(), order.Price, order.Symbol, order.Size)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(updatedOrder)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO:
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}
