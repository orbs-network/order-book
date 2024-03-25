package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/middleware"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOpenOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	logctx.Debug(r.Context(), "user trying to get their open orders", logger.String("userId", user.Id.String()))

	orders, totalOrders, err := h.svc.GetOpenOrdersForUser(r.Context(), user.Id)

	if err != nil {
		logctx.Error(r.Context(), "error getting open orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	res := middleware.NewPaginationResponse[[]models.Order](r.Context(), orders, totalOrders)

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err), logger.String("orderId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting orders. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("orderId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting orders. Try again later")
	}

}

func (h *Handler) GetFilledOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	logctx.Debug(r.Context(), "user trying to get their filled orders", logger.String("userId", user.Id.String()))

	orders, totalOrders, err := h.svc.GetFilledOrdersForUser(r.Context(), user.Id)

	if err != nil {
		logctx.Error(r.Context(), "error getting filled orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	res := middleware.NewPaginationResponse[[]models.Order](r.Context(), orders, totalOrders)

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting orders. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting orders. Try again later")
	}

}

func (h *Handler) GetSwapFills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	logctx.Debug(r.Context(), "user trying to get their filled orders", logger.String("userId", user.Id.String()))

	symbol := models.Symbol("MATIC-USDC")
	startAt, endAt := getStartEndTime(r)
	fills, err := h.svc.GetSwapFills(r.Context(), user.Id, symbol, startAt, endAt)
	if err != nil {
		logctx.Error(r.Context(), "failed GetSwapFills", logger.Error(err), logger.String("userId", user.Id.String()))
		if err == models.ErrMaxRecExceeded {
			// narrow down the time range, 256 exceeded
			restutils.WriteJSONError(ctx, w, http.StatusRequestEntityTooLarge, err.Error())
		} else {
			restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting swaps. Try again later")
		}
	}
	jsonData, err := json.Marshal(fills)
	if err != nil {
		logctx.Error(r.Context(), "failed to Marshal orders", logger.Error(err), logger.String("userId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error Marshalling swap orders.")
	}

	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("orderId", user.Id.String()))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error write response.")
	}
}
