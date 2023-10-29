package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

//	type FillOrder struct {
//		OrderID        string `json:"orderID"`
//		OrderSignatrue string `json:"orderSignatrue"`
//		AmountOut      string `json:"amountOut"`
//		Source         string `json:"source"` // 0xPubKey
//	}
type ConfirmAuctionResponse struct {
	AuctionId     string `json:"auctionId"`
	BookSignature string `json:"bookSignature"`
	//FillOrders    []FillOrder `json:"fillOrders"`
}

func (h *Handler) confirmAuction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	auctionId := vars["auctionId"]

	if auctionId == "" {
		http.Error(w, "auctionId is empty", http.StatusBadRequest)
		return
	}

	bytes := []byte(auctionId)
	uuid, err := uuid.FromBytes(bytes)
	if err != nil {
		logctx.Error(r.Context(), fmt.Sprintf("auctionID: %s", auctionId), logger.Error(err))
		logctx.Error(r.Context(), "auctionID is not a valid uuid", logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return
	}

	res, err := h.svc.ConfirmAuction(r.Context(), uuid)
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
