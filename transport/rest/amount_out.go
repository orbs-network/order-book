package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type AmountOutRequest struct {
	AuctionId string `json:"auctionId"`
	AmountIn  string `json:"amountIn"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
}

func (h *Handler) amountOut(w http.ResponseWriter, r *http.Request) {

	var args AmountOutRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	symbol, err := models.StrToSymbol(args.Symbol)
	if err != nil {
		http.Error(w, "'symbol' is not a valid", http.StatusBadRequest)
		return
	}
	amountIn, err := decimal.NewFromString(args.AmountIn)
	if err != nil {
		http.Error(w, "'size' is not a valid number format", http.StatusBadRequest)
		return
	}

	side := models.BUY
	if strings.EqualFold(args.Side, "SELL") {
		side = models.SELL
	}
	amountOutRes, err := h.svc.GetAmountOut(r.Context(), args.AuctionId, symbol, side, amountIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(amountOutRes)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal amountOutRes", logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logctx.Error(r.Context(), "failed to write amountOutRes response", logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) approveOrders(w http.ResponseWriter, r *http.Request) {
}
