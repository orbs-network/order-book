package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type ConfirmAuctionResponse struct {
	AuctionId     string `json:"auctionId"`
	BookSignature string `json:"bookSignature"`
}

type BeginAuctionReq struct {
	AmountIn string `json:"amountIn"`
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
}

type BeginAuctionRes struct {
	AuctionId string `json:"auctionId"`
	AmountOut string `json:"amountOut"`
}

func handleAuctionId(w http.ResponseWriter, r *http.Request) *uuid.UUID {
	vars := mux.Vars(r)
	auctionId := vars["auctionId"]
	ctx := r.Context()
	if auctionId == "" {
		logctx.Error(ctx, "auctionID is empty")
		http.Error(w, "auctionId is empty", http.StatusBadRequest)
		return nil
	}
	res, err := uuid.Parse(auctionId)
	if err != nil {
		logctx.Error(ctx, fmt.Sprintf("invalid auctionID: %s", auctionId), logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return nil
	}
	return &res

}

func (h *Handler) beginAuction(w http.ResponseWriter, r *http.Request) {

	auctionId := handleAuctionId(w, r)
	if auctionId == nil {
		return
	}

	var args BeginAuctionReq
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

	side, err := models.StrToSide(strings.ToLower(args.Side))
	if err != nil {
		http.Error(w, "'side' is not a valid", http.StatusBadRequest)
		return
	}

	amountOutRes, err := h.svc.GetAmountOut(r.Context(), *auctionId, symbol, side, amountIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// convert res
	baRes := BeginAuctionRes{
		AmountOut: amountOutRes.Size.String(),
		AuctionId: auctionId.String(),
	}

	resp, err := json.Marshal(baRes)
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

func (h *Handler) abortAuction(w http.ResponseWriter, r *http.Request) {
	auctionId := handleAuctionId(w, r)
	if auctionId == nil {
		return
	}
	ctx := r.Context()
	err := h.svc.RevertAuction(ctx, *auctionId)
	if err != nil {
		logctx.Error(ctx, "failed to RevertAuction", logger.Error(err))
		http.Error(w, "Error RevertAuction", http.StatusInternalServerError)
		return
	}
}

// func (h *Handler) removeAuction(w http.ResponseWriter, r *http.Request) {
// 	auctionId := handleAuctionId(w, r)
// 	if auctionId == nil {
// 		return
// 	}
// 	ctx := r.Context()
// 	err := h.svc.GetStore().RemoveAuction(ctx, *auctionId)
// 	if err != nil {
// 		logctx.Error(ctx, "failed to RemoveAuction", logger.Error(err))
// 		http.Error(w, "Error RemoveAuction", http.StatusInternalServerError)
// 		return
// 	}
// }

func (h *Handler) auctionMined(w http.ResponseWriter, r *http.Request) {

	auctionId := handleAuctionId(w, r)
	if auctionId == nil {
		return
	}
	ctx := r.Context()
	err := h.svc.AuctionMined(ctx, *auctionId)
	if err != nil {
		logctx.Error(ctx, "failed to AuctionMined", logger.Error(err))
		http.Error(w, "Error AuctionMined", http.StatusInternalServerError)
		return
	}

}
