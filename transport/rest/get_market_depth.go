package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
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
	vars := mux.Vars(r)
	symbolStr, ok := vars["symbol"]
	if !ok {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	symbol, err := models.StrToSymbol(symbolStr)
	if err != nil {
		http.Error(w, "Invalid symbol", http.StatusBadRequest)
		return
	}

	// Getting the limit query parameter, if not provided set to default
	limitStr := r.URL.Query().Get("limit")
	limit := DEFAULT_LIMIT
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > MAX_LIMIT {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	if limit <= 0 || limit > MAX_LIMIT {
		http.Error(w, fmt.Sprintf("Invalid limit: must be between 1 and %d", MAX_LIMIT), http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to get market depth", logger.String("symbol", symbol.String()), logger.Int("limit", limit))
	marketDepth, err := h.svc.GetMarketDepth(r.Context(), symbol, limit)

	if err != nil {
		http.Error(w, "Error getting market depth", http.StatusInternalServerError)
		return
	}

	res := MarketDepthResponse{
		Code: "OK",
		Data: marketDepth,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal market depth", logger.Error(err))
		http.Error(w, "Error getting market depth", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		http.Error(w, "Error getting market depth", http.StatusInternalServerError)
	}
}
