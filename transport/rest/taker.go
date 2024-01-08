package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type BeginSwapReq struct {
	AmountIn string `json:"amountIn"`
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
}

type BeginSwapRes struct {
	SwapId    string `json:"swapId"`
	AmountOut string `json:"amountOut"`
}

type Fragment struct {
	OrderId       string                 `json:"orderId"`
	AmountOut     string                 `json:"amountOut"`
	Eip712Sig     string                 `json:"eip712Sig"`
	Eip712MsgData map[string]interface{} `json:"eip712MsgData"`
}
type ConfirmSwapRes struct {
	SwapId        string     `json:"swapId"`
	Fragments     []Fragment `json:"fragments"`
	BookSignature string     `json:"bookSignature"`
}

type QuoteReq struct {
	InAmount string `json:"inAmount"`
	InToken  string `json:"inToken"`
	OutToken string `json:"outToken"`
}

type QuoteRes struct {
	OutAmount string     `json:"outAmount"`
	OutToken  string     `json:"outToken"`
	InAmount  string     `json:"inAmount"`
	InToken   string     `json:"inToken"`
	SwapId    string     `json:"swapId"`
	Fragments []Fragment `json:"fragments"`
	//BookSignature? string     `json:"bookSignature"`
}

func (h *Handler) convertToTokenDec(ctx context.Context, outToken string, outAmount decimal.Decimal) string {
	if token, ok := h.supportedTokens[strings.ToUpper(outToken)]; ok {
		decMul := math.Pow10(token.Decimals)
		mul := outAmount.Mul(decimal.NewFromInt(int64(decMul)))
		return mul.Truncate(0).String()
	}
	logctx.Error(ctx, "Token is not found in supported tokens: "+outToken)
	return ""

}
func (h *Handler) handleQuote(w http.ResponseWriter, r *http.Request, isSwap bool) {
	var req QuoteReq
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logctx.Error(ctx, "handleQuote - failed to decode body", logger.Error(err))
		http.Error(w, "Quote invalid JSON body", http.StatusBadRequest)
		return
	}

	inAmount, err := decimal.NewFromString(req.InAmount)
	if err != nil {
		logctx.Error(ctx, "'Quote::inAmount' is not a valid number format", logger.Error(err))
		http.Error(w, "'Quote::inAmount' is not a valid number format", http.StatusBadRequest)
		return
	}

	pair := h.pairMngr.Resolve(req.InToken, req.OutToken)
	if pair == nil {
		msg := fmt.Sprintf("no suppoerted pair with found with the following tokens %s, %s", req.InToken, req.OutToken)
		logctx.Error(ctx, msg, logger.Error(err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	side := pair.GetSide(req.InToken)
	amountOutRes, err := h.svc.GetQuote(r.Context(), pair.Symbol(), side, inAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	convOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, amountOutRes.Size)
	if convOutAmount == "" {
		http.Error(w, models.ErrTokenNotsupported.Error(), http.StatusBadRequest)
		return
	}
	// convert res
	quoteRes := QuoteRes{
		OutAmount: convOutAmount,
		OutToken:  req.OutToken,
		InAmount:  req.InAmount,
		InToken:   req.InToken,
		SwapId:    "",
		Fragments: []Fragment{},
	}

	if isSwap {
		swapData, err := h.svc.BeginSwap(r.Context(), amountOutRes) //, pair.Symbol())//, side, inAmount)
		if err != nil {
			http.Error(w, "BeginSwap filed", http.StatusBadRequest)
			return
		}
		for i := 0; i < len(swapData.Fragments); i++ {
			convOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, swapData.Fragments[i].Size)
			frag := Fragment{
				OrderId:       swapData.Fragments[i].OrderId.String(),
				AmountOut:     convOutAmount,
				Eip712Sig:     swapData.Orders[i].Signature.Eip712Sig,
				Eip712MsgData: swapData.Orders[i].Signature.Eip712MsgData,
			}
			quoteRes.Fragments = append(quoteRes.Fragments, frag)
		}
		quoteRes.SwapId = swapData.SwapId.String()
	}

	resp, err := json.Marshal(quoteRes)
	if err != nil {
		logctx.Error(r.Context(), "failed to marshal QuoteRes", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		logctx.Error(r.Context(), "failed to write QuoteRes response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Quote METHOD GET
func (h *Handler) quote(w http.ResponseWriter, r *http.Request) {
	h.handleQuote(w, r, false)
}

// SWAP METHOD GET
func (h *Handler) swap(w http.ResponseWriter, r *http.Request) {
	h.handleQuote(w, r, true)
}

// Helper
func handleSwapId(w http.ResponseWriter, r *http.Request) *uuid.UUID {
	vars := mux.Vars(r)
	swapId := vars["swapId"]
	ctx := r.Context()
	if swapId == "" {
		logctx.Error(ctx, "swapID is empty")
		http.Error(w, "swapId is empty", http.StatusBadRequest)
		return nil
	}
	res, err := uuid.Parse(swapId)
	if err != nil {
		logctx.Error(ctx, fmt.Sprintf("invalid swapID: %s", swapId), logger.Error(err))
		http.Error(w, "Error GetAmountOut", http.StatusInternalServerError)
		return nil
	}
	return &res
}

// POST
func (h *Handler) abortSwap(w http.ResponseWriter, r *http.Request) {
	swapId := handleSwapId(w, r)
	if swapId == nil {
		return
	}
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	if err := h.svc.AbortSwap(r.Context(), *swapId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Create an empty JSON object
		obj := genRes{
			StatusText: err.Error(),
			Status:     http.StatusBadRequest,
		}

		// Convert the emptyJSON object to JSON format
		jRes, _ := json.Marshal(obj)
		_, err = w.Write(jRes)
		if err != nil {
			logctx.Error(r.Context(), "abortSwap - failed to write resp", logger.Error(err))
		}
	}

	// Write the JSON response with a status code of 200
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(h.okJson)
	if err != nil {
		logctx.Error(r.Context(), "abortSwap - failed to write resp", logger.Error(err))
	}

}
