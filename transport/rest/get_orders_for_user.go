package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/middleware"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOpenOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	logctx.Info(r.Context(), "user trying to get their open orders", logger.String("userId", user.Id.String()))

	orders, totalOrders, err := h.svc.GetOpenOrdersForUser(r.Context(), user.Id)

	if err != nil {
		logctx.Error(r.Context(), "error getting open orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := middleware.NewPaginationResponse[[]models.Order](r.Context(), orders, totalOrders)

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err), logger.String("orderId", user.Id.String()))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("orderId", user.Id.String()))
		http.Error(w, "Error getting orders. Try again later", http.StatusInternalServerError)
	}

}

func (h *Handler) GetFilledOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	logctx.Info(r.Context(), "user trying to get their filled orders", logger.String("userId", user.Id.String()))

	orders, totalOrders, err := h.svc.GetFilledOrdersForUser(r.Context(), user.Id)

	if err != nil {
		logctx.Error(r.Context(), "error getting filled orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := middleware.NewPaginationResponse[[]models.Order](r.Context(), orders, totalOrders)

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err), logger.String("orderId", user.Id.String()))
		http.Error(w, "Error creating order. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err), logger.String("orderId", user.Id.String()))
		http.Error(w, "Error getting orders. Try again later", http.StatusInternalServerError)
	}

}
