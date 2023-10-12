package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrderResponse struct {
	OrderId string `json:"orderId"`
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIdStr := vars["orderId"]

	orderId, err := uuid.Parse(orderIdStr)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to cancel order", logger.String("userId", userId.String()), logger.String("orderId", orderId.String()))
	err = h.svc.CancelOrder(r.Context(), orderId)

	if err == models.ErrOrderNotFound {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err != nil {
		logctx.Error(r.Context(), "failed to cancel order", logger.Error(err))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	res := CancelOrderResponse{
		OrderId: orderId.String(),
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal created order", logger.Error(err))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
