package rest

import (
	"context"
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

func (h *Handler) CancelOrderByOrderId(w http.ResponseWriter, r *http.Request) {
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

	logctx.Info(userCtx, "user trying to cancel order by orderID", logger.String("userId", userId.String()), logger.String("orderId", orderId.String()))

	h.handleCancelOrder(userCtx, orderId, false, w)
}

func (h *Handler) CancelOrderByClientOId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientOIdStr := vars["clientOId"]

	clientOId, err := uuid.Parse(clientOIdStr)
	if err != nil {
		http.Error(w, "invalid clientOId", http.StatusBadRequest)
		return
	}

	// TODO: don't hardcode user
	user := models.User{
		ID:   userId,
		Type: models.MARKET_MAKER,
	}

	userCtx := utils.WithUser(r.Context(), &user)

	logctx.Info(userCtx, "user trying to cancel order by clientOId", logger.String("userId", userId.String()), logger.String("clientOId", clientOId.String()))

	h.handleCancelOrder(userCtx, clientOId, true, w)
}

// handleCancelOrder calls the service to cancel an order and writes the response to the client
func (h *Handler) handleCancelOrder(userCtx context.Context, id uuid.UUID, isClientOId bool, w http.ResponseWriter) {

	cancelledOrderId, err := h.svc.CancelOrder(userCtx, id, isClientOId)

	if err == models.ErrNoUserInContext {
		logctx.Error(userCtx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err == models.ErrOrderNotFound {
		logctx.Warn(userCtx, "order not found", logger.String("id", id.String()))
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err == models.ErrOrderNotOpen {
		logctx.Warn(userCtx, "user trying to cancel order that is not open", logger.String("id", id.String()))
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err == models.ErrUnauthorized {
		logctx.Warn(userCtx, "user not authorized to cancel order", logger.String("id", id.String()))
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		logctx.Error(userCtx, "failed to cancel order", logger.Error(err))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	if cancelledOrderId == nil {
		logctx.Error(userCtx, "cancelled order ID is nil")
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	res := CancelOrderResponse{
		OrderId: cancelledOrderId.String(),
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
		logctx.Error(userCtx, "failed to write response", logger.Error(err), logger.String("orderId", cancelledOrderId.String()))
		http.Error(w, "Error cancelling order. Try again later", http.StatusInternalServerError)
	}
}
