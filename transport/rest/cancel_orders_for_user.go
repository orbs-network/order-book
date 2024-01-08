package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
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
		restutils.WriteJSONError(w, http.StatusUnauthorized, "User not found")
		return
	}

	logctx.Info(ctx, "user trying to cancel all their orders", logger.String("userId", user.Id.String()))
	orderIds, err := h.svc.CancelOrdersForUser(ctx, user.Id)

	if err == models.ErrNotFound {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(w, http.StatusNotFound, "No orders found")
		return
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Unable to cancel orders. Try again later")
		return
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", user.Id.String()), logger.Int("numOrders", len(orderIds)))

	res := CancelOrdersForUserResponse{
		CancelledOrderIds: orderIds,
	}

	orderIdsJSON, err := json.Marshal(res)
	if err != nil {
		logctx.Error(ctx, "could not marshal orderIds to JSON", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Unable to marshal orderIds to JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(orderIdsJSON); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error cancelling orders. Try again later")
	}
}
