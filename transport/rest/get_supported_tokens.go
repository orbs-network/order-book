package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

var filePath = os.Getenv("SUPPORTED_TOKENS_JSON_FILE_PATH")

type res struct {
	Tokens service.SupportedTokens `json:"tokens"`
}

func (h *Handler) GetSupportedTokens(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if filePath == "" {
		logctx.Warn(r.Context(), "SUPPORTED_TOKENS_JSON_FILE_PATH env var not set, using default")
		filePath = "supportedTokens.json"
	}

	logctx.Info(ctx, "User requesting supported tokens", logger.String("user", user.Id.String()))
	tokens, err := service.GetSupportedTokens(r.Context(), filePath)
	if err != nil {
		logctx.Error(r.Context(), "failed to get supported tokens", logger.Error(err))
		http.Error(w, "Error getting supported tokens", http.StatusInternalServerError)
		return
	}

	res := res{Tokens: tokens}

	jsonData, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal response", logger.Error(err))
		http.Error(w, "Error getting supported tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		http.Error(w, "Error getting supported tokens", http.StatusInternalServerError)
	}
}
