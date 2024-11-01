package rest

import (
	"net/http"

	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetSupportedTokens(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	jsonData, err := h.supportedTokens.AsJson()

	if err != nil {
		logctx.Error(ctx, "failed to marshal response", logger.Error(err))
		http.Error(w, "Error getting supported tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(ctx, "failed to write response", logger.Error(err))
		http.Error(w, "Error getting supported tokens", http.StatusInternalServerError)
	}
}
