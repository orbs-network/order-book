package rest

import (
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

func (h *Handler) amountOut(w http.ResponseWriter, r *http.Request) {

	var args CreateOrderRequest
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
	decSize, err := decimal.NewFromString(args.Size)
	if err != nil {
		http.Error(w, "'size' is not a valid number format", http.StatusBadRequest)
		return
	}

	//side = getSide
	h.svc.GetAmountOut(nil, true, symbol, decSize)

}

func (h *Handler) approveOrders(w http.ResponseWriter, r *http.Request) {
}
