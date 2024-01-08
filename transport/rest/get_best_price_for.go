package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type GetBestPriceForResponse struct {
	Price  string        `json:"price"`
	Side   models.Side   `json:"side"`
	Symbol models.Symbol `json:"symbol"`
}

func (h *Handler) GetBestPriceFor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(w, http.StatusUnauthorized, "User not found")
		return
	}

	vars := mux.Vars(r)
	symbolStr := vars["symbol"]
	sideStr := vars["side"]

	symbol, err := models.StrToSymbol(symbolStr)
	if err != nil {
		restutils.WriteJSONError(w, http.StatusBadRequest, "invalid symbol")
		return
	}

	side, err := models.StrToSide(sideStr)
	if err != nil {
		restutils.WriteJSONError(w, http.StatusBadRequest, "invalid side")
		return
	}

	logctx.Info(r.Context(), "user trying to get best price", logger.String("userId", user.Id.String()), logger.String("symbol", symbol.String()), logger.String("side", side.String()))
	price, err := h.svc.GetBestPriceFor(r.Context(), symbol, side)

	if err != nil {
		logctx.Error(r.Context(), "failed to get best price", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting best price. Try again later")
		return
	}

	if price.IsZero() {
		logctx.Info(r.Context(), "no orders found for symbol and side", logger.String("symbol", symbol.String()), logger.String("side", side.String()))
		restutils.WriteJSONError(w, http.StatusNotFound, "No orders found for symbol and side")
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
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting best price. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		restutils.WriteJSONError(w, http.StatusInternalServerError, "Error getting best price. Try again later")
	}
}
