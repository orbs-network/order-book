package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type BeginSwapRes struct {
	SwapId   string `json:"swapId"`
	QuoteRes string `json:"quoteRes"`
}

type Fragment struct {
	Signature string             `json:"signature"`
	Abi       string             `json:"string"`
	AbiData   models.AbiFragment `json:"abiData"`
}
type ConfirmSwapRes struct {
	SwapId        string     `json:"swapId"`
	Fragments     []Fragment `json:"fragments"`
	BookSignature string     `json:"bookSignature"`
}

type QuoteReq struct {
	InAmount     string `json:"inAmount"`
	InToken      string `json:"inToken"`
	OutToken     string `json:"outToken"`
	MinOutAmount string `json:"minOutAmount"`
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

func (h *Handler) convertToTokenDec(ctx context.Context, tokenName string, amount decimal.Decimal) string {
	if token := h.supportedTokens.ByName(tokenName); token != nil {
		decMul := math.Pow10(token.Decimals)
		mul := amount.Mul(decimal.NewFromInt(int64(decMul)))
		return mul.Truncate(0).String()
	}
	logctx.Error(ctx, "Token is not found in supported tokens: "+tokenName)
	return ""
}

func (h *Handler) convertFromTokenDec(ctx context.Context, tokenName, amountStr string) (decimal.Decimal, error) {
	if token := h.supportedTokens.ByName(tokenName); token != nil {
		decDiv := math.Pow10(token.Decimals)
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			logctx.Error(ctx, "error converting amountStr: "+amountStr)
			return decimal.Zero, err
		}
		res := amount.Div(decimal.NewFromInt(int64(decDiv)))
		return res, nil
	}
	logctx.Error(ctx, "Token is not found in supported tokens: "+tokenName)
	return decimal.Zero, models.ErrTokenNotsupported
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

	inAmount, err := h.convertFromTokenDec(ctx, req.InToken, req.InAmount)
	if err != nil {
		logctx.Error(ctx, "'Quote::inAmount' is not a valid number format", logger.Error(err))
		http.Error(w, "'Quote::inAmount' is not a valid number format", http.StatusBadRequest)
		return
	}

	// a threshold for min amount out expect, return error if
	var minOutAmount *decimal.Decimal = nil
	if req.MinOutAmount != "" {
		convMinOutAmount, err := h.convertFromTokenDec(ctx, req.OutToken, req.MinOutAmount)
		if err != nil {
			logctx.Warn(ctx, "'Quote::minOutAmount' is not a valid number format - passing nil", logger.Error(err))
		} else {
			minOutAmount = &convMinOutAmount
		}
	}

	pair := h.pairMngr.Resolve(req.InToken, req.OutToken)
	if pair == nil {
		msg := fmt.Sprintf("no suppoerted pair with found with the following tokens %s, %s", req.InToken, req.OutToken)
		logctx.Error(ctx, msg, logger.Error(err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	side := pair.GetSide(req.InToken)
	svcQuoteRes, err := h.svc.GetQuote(r.Context(), pair.Symbol(), side, inAmount, minOutAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	convOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, svcQuoteRes.Size)
	if convOutAmount == "" {
		http.Error(w, models.ErrTokenNotsupported.Error(), http.StatusBadRequest)
		return
	}
	// convert res
	res := QuoteRes{
		OutAmount: convOutAmount,
		OutToken:  req.OutToken,
		InAmount:  req.InAmount,
		InToken:   req.InToken,
		SwapId:    "",
		Fragments: []Fragment{},
	}

	if isSwap {
		swapData, err := h.svc.BeginSwap(r.Context(), svcQuoteRes)
		if err != nil {
			http.Error(w, "BeginSwap filed", http.StatusBadRequest)
			return
		}
		// inToken, ok := h.supportedTokens[req.InToken]
		// if !ok {
		// 	http.Error(w, "InToken address not found", http.StatusBadRequest)
		// 	return
		// }
		// outToken, ok := h.supportedTokens[req.OutToken]
		// if !ok {
		// 	http.Error(w, "InToken address not found", http.StatusBadRequest)
		// 	return
		// }

		for i := 0; i < len(swapData.Fragments); i++ {
			// conver In/Out amount to token decimals
			convInAmount := h.convertToTokenDec(r.Context(), req.InToken, swapData.Fragments[i].InSize)
			convOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, swapData.Fragments[i].OutSize)

			// convert to uint256 for abi encode
			inputAmount := big.NewInt(0)
			inputAmount.SetString(convInAmount, 10)

			outputAmount := big.NewInt(0)
			outputAmount.SetString(convOutAmount, 10)

			abiData := swapData.Orders[i].Signature.AbiFragment
			abiData.ExclusivityOverrideBps = big.NewInt(0)
			abiData.Input.Amount = inputAmount
			if len(abiData.Outputs) == 0 {
				logctx.Error(r.Context(), "abiData.Outputs length is 0")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			abiData.Outputs[0].Amount.SetString(convOutAmount, 10)
			abiEncoded, err := models.EncodeFragData(ctx, abiData)

			if err != nil {
				logctx.Error(ctx, "args.Pack failed %s", logger.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			frag := Fragment{
				Signature: swapData.Orders[i].Signature.Eip712Sig,
				Abi:       abiEncoded,
				AbiData:   abiData,
			}
			res.Fragments = append(res.Fragments, frag)
		}
		res.SwapId = swapData.SwapId.String()
	}

	resp, err := json.Marshal(res)
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

// Quote METHOD POST
func (h *Handler) quote(w http.ResponseWriter, r *http.Request) {
	h.handleQuote(w, r, false)
}

// SWAP METHOD POST
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
		http.Error(w, "Error GetQuoteRes", http.StatusInternalServerError)
		return nil
	}
	return &res
}

// POST
func (h *Handler) swapStarted(w http.ResponseWriter, r *http.Request) {
	swapId := handleSwapId(w, r)
	if swapId == nil {
		return
	}

	// get txHash
	vars := mux.Vars(r)
	txhash := vars["txHash"]
	ctx := r.Context()
	if txhash == "" {
		logctx.Error(ctx, "txhash is empty")
		http.Error(w, "txhash is empty", http.StatusBadRequest)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	if err := h.svc.SwapStarted(r.Context(), *swapId, txhash); err != nil {
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
