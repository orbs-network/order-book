package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIdStr := vars["orderId"]

	orderId, err := uuid.Parse(orderIdStr)
	if err != nil {
		http.Error(w, "invalid orderId", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to get order by ID", logger.String("orderId", orderId.String()))
	order, err := h.svc.GetOrderById(r.Context(), orderId)

	if err != nil {
		http.Error(w, "Internal error. Try again later", http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, fmt.Sprintf("order not found for %q", orderIdStr), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(order)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal order", logger.Error(err))
		http.Error(w, "Error getting order price", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		http.Error(w, "Error getting order by ID", http.StatusInternalServerError)
	}
}
