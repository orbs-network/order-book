package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	vars := mux.Vars(r)
	orderIdStr := vars["orderId"]

	orderId, err := uuid.Parse(orderIdStr)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Invalid orderId")
		return
	}

	logctx.Info(r.Context(), "user trying to get order by ID", logger.String("userId", user.Id.String()), logger.String("orderId", orderId.String()))
	order, err := h.svc.GetOrderById(r.Context(), orderId)

	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Internal error. Try again later")
		return
	}

	if order == nil {
		restutils.WriteJSONError(ctx, w, http.StatusNotFound, fmt.Sprintf("Order not found for %s", orderIdStr))
		return
	}

	resp, err := json.Marshal(order)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal order", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting order by ID")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting order by ID")
	}
}
