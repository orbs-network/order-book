package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
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

	// TODO: don't hardcode user
	user := models.User{
		ID:   userId,
		Type: models.MARKET_MAKER,
	}

	userCtx := utils.WithUser(r.Context(), &user)

	logctx.Info(userCtx, "user trying to cancel order", logger.String("userId", userId.String()), logger.String("orderId", orderId.String()))
	err = h.svc.CancelOrder(userCtx, orderId)

	if err == models.ErrNoUserInContext {
		logctx.Error(userCtx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err == models.ErrOrderNotFound {
		logctx.Warn(userCtx, "order not found", logger.String("orderId", orderId.String()))
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err == models.ErrOrderNotOpen {
		logctx.Warn(userCtx, "user trying to cancel order that is not open", logger.String("orderId", orderId.String()))
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err == models.ErrUnauthorized {
		logctx.Warn(userCtx, "user not authorized to cancel order", logger.String("orderId", orderId.String()))
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		logctx.Error(userCtx, "failed to cancel order", logger.Error(err))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	res := CancelOrderResponse{
		OrderId: orderId.String(),
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(userCtx, "failed to marshal created order", logger.Error(err))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(userCtx, "failed to write response", logger.Error(err), logger.String("orderId", orderId.String()))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
	}

}
