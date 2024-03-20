package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrdersForUserResponse struct {
	Symbol            string      `json:"symbol"`
	CancelledOrderIds []uuid.UUID `json:"cancelledOrderIds"`
}

func (h *Handler) CancelOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	// get symbol
	symbol := r.URL.Query().Get("symbol")
	symbol = strings.ToUpper(symbol)
	if symbol == "" {
		logctx.Info(ctx, "cancelAll symbol was not provided, cancelling all orders in all symbols", logger.String("userId", user.Id.String()))
	}

	logctx.Info(ctx, "user trying to cancel all their orders", logger.String("symbol", symbol), logger.String("userId", user.Id.String()))
	orderIds, err := h.svc.CancelOrdersForUser(ctx, user.Id, models.Symbol(symbol))

	if err == models.ErrNotFound {
		logctx.Info(ctx, "no orders found for user", logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusNotFound, "No orders found")
		return
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Unable to cancel orders. Try again later")
		return
	}

	logctx.Info(ctx, "cancelled all orders for user", logger.String("userId", user.Id.String()), logger.Int("numOrders", len(orderIds)))

	res := CancelOrdersForUserResponse{
		Symbol:            symbol,
		CancelledOrderIds: orderIds,
	}

	orderIdsJSON, err := json.Marshal(res)
	if err != nil {
		logctx.Error(ctx, "could not marshal orderIds to JSON", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Unable to marshal orderIds to JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(orderIdsJSON); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error cancelling orders. Try again later")
	}
}
