package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type MarketDepthResponse struct {
	Code string             `json:"code"`
	Data models.MarketDepth `json:"data"`
}

const MAX_LIMIT int = 1000
const DEFAULT_LIMIT int = 10

func (h *Handler) GetMarketDepth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	vars := mux.Vars(r)
	symbolStr, ok := vars["symbol"]
	if !ok {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Symbol is required")
		return
	}

	symbol, err := models.StrToSymbol(symbolStr)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Invalid symbol")
		return
	}

	// Getting the limit query parameter, if not provided set to default
	limitStr := r.URL.Query().Get("limit")
	limit := DEFAULT_LIMIT
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > MAX_LIMIT {
			restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Invalid limit")
			return
		}
	}

	if limit <= 0 || limit > MAX_LIMIT {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("Invalid limit: must be between 1 and %d", MAX_LIMIT))
		return
	}

	logctx.Info(r.Context(), "user trying to get market depth", logger.String("userId", user.Id.String()), logger.String("symbol", symbol.String()), logger.Int("limit", limit))
	marketDepth, err := h.svc.GetMarketDepth(r.Context(), symbol, limit)

	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting market depth. Try again later")
		return
	}

	res := MarketDepthResponse{
		Code: "OK",
		Data: marketDepth,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal market depth", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting market depth. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "Error getting market depth. Try again later")
	}
}
