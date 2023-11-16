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
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	logctx.Info(ctx, "user trying to cancel all their orders", logger.String("userId", user.Id.String()))
	err := h.svc.CancelOrdersForUser(ctx, user.Id)

	if err == models.ErrUserNotFound {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		logctx.Error(ctx, "could not cancel orders for user", logger.Error(err), logger.String("userId", user.Id.String()))
		http.Error(w, "Unable to cancel orders. Try again later", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
