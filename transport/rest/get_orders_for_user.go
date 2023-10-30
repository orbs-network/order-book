package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOrdersForUser(w http.ResponseWriter, r *http.Request) {
	logctx.Info(r.Context(), "user trying to get their orders", logger.String("user_id", userId.String()))

	orders, totalOrders, err := h.svc.GetOrdersForUser(r.Context(), userId)

	if err != nil {
		logctx.Error(r.Context(), "error getting orders for user", logger.Error(err), logger.String("user_id", userId.String()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := NewPaginationResponse[[]models.Order](r.Context(), orders, totalOrders)

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err), logger.String("orderId", userId.String()))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("orderId", userId.String()))
		http.Error(w, "Error getting orders. Try again later", http.StatusInternalServerError)
	}

}
