package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrdersForUserResponse struct {
	CancelledOrderIds []uuid.UUID `json:"cancelledOrderIds"`
}

func (h *Handler) CancelOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	logctx.Info(ctx, "user trying to cancel all their orders", logger.String("userId", user.Id.String()))
	orderIds, err := h.svc.CancelOrdersForUser(ctx, user.Id)

	if err == models.ErrNoOrdersFound {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", user.Id.String()))
		http.Error(w, "No orders found", http.StatusNotFound)
		return
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, "Unable to cancel orders. Try again later", http.StatusInternalServerError)
		return
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", user.Id.String()), logger.Int("numOrders", len(orderIds)))

	res := CancelOrdersForUserResponse{
		CancelledOrderIds: orderIds,
	}

	orderIdsJSON, err := json.Marshal(res)
	if err != nil {
		logctx.Error(ctx, "could not marshal orderIds to JSON", logger.Error(err))
		http.Error(w, "Unable to marshal orderIds to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(orderIdsJSON); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, "Error cancelling orders. Try again later", http.StatusInternalServerError)
	}
}
