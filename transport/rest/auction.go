package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type FillOrder struct {
	OrderID        string `json:"orderID"`
	OrderSignatrue string `json:"orderSignatrue"`
	AmountOut      string `json:"amountOut"`
}
type ConfirmAuctionResponse struct {
	AuctionId     string      `json:"auctionId"`
	BookSignature string      `json:"bookSignature"`
	FillOrders    []FillOrder `json:"fillOrders"`
}

func (h *Handler) confirmAuction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	auctionId := vars["auctionId"]

	if auctionId == "" {
		http.Error(w, "auctionId is empty", http.StatusBadRequest)
		return
	}

	res, err := h.svc.ConfirmAuction(r.Context(), auctionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal confirmAuction response", logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logctx.Error(r.Context(), "failed to write confirmAuction response", logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) removeAuction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	auctionId := vars["auctionId"]
	if auctionId == "" {
		http.Error(w, "auctionId is empty", http.StatusBadRequest)
		return
	}

}
