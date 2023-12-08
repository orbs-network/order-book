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

type BeginAuctionReq struct {
	AmountIn string `json:"amountIn"`
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
}

type BeginAuctionRes struct {
	AuctionId string `json:"auctionId"`
	AmountOut string `json:"amountOut"`
}

type Fragment struct {
	OrderId       string                 `json:"orderId"`
	AmountOut     string                 `json:"amountOut"`
	Eip712Sig     string                 `json:"eip712Sig"`
	Eip712MsgData map[string]interface{} `json:"eip712MsgData"`
}
type ConfirmAuctionRes struct {
	AuctionId     string     `json:"auctionId"`
	Fragments     []Fragment `json:"fragments"`
	BookSignature string     `json:"bookSignature"`
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logctx.Error(r.Context(), "failed to write amountOutRes response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Handler) confirmAuction(w http.ResponseWriter, r *http.Request) {

	auctionId := handleAuctionId(w, r)
	if auctionId == nil {
		return
	}

	caRes, err := h.svc.ConfirmAuction(r.Context(), *auctionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert service response to rest response
	respObj := ConfirmAuctionRes{
		AuctionId:     auctionId.String(),
		BookSignature: "TODO:sign",
	}
	for i := 0; i < len(caRes.Fragments); i++ {
		frag := Fragment{
			AmountOut:     caRes.Fragments[i].Size.String(),
			OrderId:       caRes.Orders[i].Id.String(),
			Eip712Sig:     caRes.Orders[i].Signature.Eip712Sig,
			Eip712MsgData: caRes.Orders[i].Signature.Eip712MsgData,
		}
		respObj.Fragments = append(respObj.Fragments, frag)
	}
	resp, err := json.Marshal(respObj)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal confirmAuction response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logctx.Error(r.Context(), "failed to write confirmAuction response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) auctionMined(w http.ResponseWriter, r *http.Request) {
	auctionId := handleAuctionId(w, r)
	if auctionId == nil {
		return
	}
	ctx := r.Context()
	err := h.svc.AuctionMined(ctx, *auctionId)
	if err != nil {
		logctx.Error(ctx, "failed to AuctionMined", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
