package rest

import (
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) CancelOrdersForUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	publicKey := utils.GetPkCtx(ctx)
	if publicKey == "" {
		http.Error(w, "Missing public key", http.StatusBadRequest)
		return
	}

	logctx.Info(ctx, "user trying to cancel all their orders", logger.String("publicKey", publicKey))
	err := h.svc.CancelOrdersForUser(ctx, publicKey)

	if err == models.ErrUserNotFound {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("publicKey", publicKey))
		http.Error(w, "Unable to cancel orders. Try again later", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}