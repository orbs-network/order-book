package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type symbol struct {
	Symbol string `json:"symbol"`
	// TODO: add more fields
}

type getSymbolsResponse []symbol

func (h *Handler) GetSymbols(w http.ResponseWriter, r *http.Request) {
	symbols, err := h.svc.GetSymbols(r.Context())
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal symbols", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting symbols. Try again later")
		return
	}

	symbolsSlice := getSymbolsResponse{}
	for _, s := range symbols {
		symbolsSlice = append(symbolsSlice, symbol{Symbol: s.String()})
	}

	resp, err := json.Marshal(symbolsSlice)

	if err != nil {
		logctx.Error(r.Context(), "failed to marshal symbols", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting symbols. Try again later")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting symbols. Try again later")
	}
}
