package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type GetBestPriceForResponse struct {
	Price  string        `json:"price"`
	Side   models.Side   `json:"side"`
	Symbol models.Symbol `json:"symbol"`
}

func (h *Handler) GetBestPriceFor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbolStr := vars["symbol"]
	sideStr := vars["side"]

	symbol, err := models.StrToSymbol(symbolStr)
	if err != nil {
		http.Error(w, "invalid symbol", http.StatusBadRequest)
		return
	}

	side, err := models.StrToSide(sideStr)
	if err != nil {
		http.Error(w, "invalid side", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to get best price", logger.String("userId", userId.String()), logger.String("symbol", symbol.String()), logger.String("side", side.String()))
	price, err := h.svc.GetBestPriceFor(r.Context(), symbol, side)

	if err != nil {
		logctx.Error(r.Context(), "failed to get best price", logger.Error(err))
		http.Error(w, "Error getting best price. Try again later", http.StatusInternalServerError)
		return
	}

	if price.IsZero() {
		logctx.Info(r.Context(), "no orders found for symbol and side", logger.String("symbol", symbol.String()), logger.String("side", side.String()))
		http.Error(w, "No orders found for symbol and side", http.StatusNotFound)
		return
	}

	res := GetBestPriceForResponse{
		Price:  price.String(),
		Side:   side,
		Symbol: symbol,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal best price", logger.Error(err))
		http.Error(w, "Error getting best price. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		http.Error(w, "Error getting best price. Try again later", http.StatusInternalServerError)
	}
}
