package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/abi"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

type BeginSwapRes struct {
	SwapId   string `json:"swapId"`
	QuoteRes string `json:"quoteRes"`
}

type Fragment struct {
	Signature      string    `json:"signature"`
	AbiOrder       abi.Order `json:"abiOrder"`
	TakerInAmount  string    `json:"takerInAmount"`
	TakerOutAmount string    `json:"takerOutAmount"`
}

type QuoteReq struct {
	InAmount        string `json:"inAmount"`
	InToken         string `json:"inToken"`
	InTokenAddress  string `json:"inTokenAddress"`
	OutToken        string `json:"outToken"`
	OutTokenAddress string `json:"outTokenAddress"`
	MinOutAmount    string `json:"minOutAmount"`
}

type QuoteRes struct {
	OutAmount string     `json:"outAmount"`
	OutToken  string     `json:"outToken"`
	InAmount  string     `json:"inAmount"`
	InToken   string     `json:"inToken"`
	SwapId    string     `json:"swapId"`
	AbiCall   string     `json:"abiCall"`
	Contract  string     `json:"contract"`
	Fragments []Fragment `json:"fragments"`
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

// returns resolve name
// only if
// 1. name is missing
// 2. address exists
// 3. address is found in supported tokens
// returns error if needed
// returns empty string if no need to resolve
func (h *Handler) nameFromAddress(name, address string) (string, error) {
	token := h.supportedTokens.ByAddress(address)
	if token == nil {
		return "", models.ErrTokenNotsupported
	}
	return token.Name, nil
}

func (h *Handler) resolveQuoteTokenNames(req *QuoteReq) error {
	// has address but no name
	if req.InToken == "" {
		InName, err := h.nameFromAddress(req.InToken, req.InTokenAddress)
		if err != nil {
			return err
		}
		if len(InName) > 0 {
			req.InToken = InName
		}
	}
	if req.OutToken == "" {
		OutName, err := h.nameFromAddress(req.OutToken, req.OutTokenAddress)
		if err != nil {
			return err
		}
		if len(OutName) > 0 {
			req.OutToken = OutName
		}
	}
	return nil
}
func Signature2Bytes(sig string) []byte {
	// remove leading 0x if exists
	sig = strings.TrimPrefix(sig, "0x")
	return common.Hex2Bytes(sig)
}

func (h *Handler) handleQuote(w http.ResponseWriter, r *http.Request, isSwap bool) *QuoteRes {
	var req QuoteReq
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error())
		return nil
	}

	// ensure token names if only addresses were sent
	err = h.resolveQuoteTokenNames(&req)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error(), logger.String("InTokenAddress", req.InTokenAddress), logger.String("OutTokenAddress", req.OutTokenAddress))
		return nil
	}

	inAmount, err := h.convertFromTokenDec(ctx, req.InToken, req.InAmount)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error(), logger.String("InToken", req.InToken), logger.Error(models.ErrTokenNotsupported))
		return nil
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
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "no suppoerted pair was found for tokens", logger.String("InToken", req.InToken), logger.String("OutToken", req.OutToken))
		return nil
	}
	side := pair.GetSide(req.InToken)
	svcQuoteRes, err := h.svc.GetQuote(r.Context(), pair.Symbol(), side, inAmount, minOutAmount)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
		return nil
	}

	convOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, svcQuoteRes.Size)
	if convOutAmount == "" {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error(), logger.String("OutToken", req.OutToken))
		return nil
	}
	// convert res
	res := QuoteRes{
		OutAmount: convOutAmount,
		OutToken:  req.OutToken,
		InAmount:  req.InAmount,
		InToken:   req.InToken,
		//SwapId:    "",
		Fragments: []Fragment{},
	}

	// debug
	//AbiElements

	if isSwap {
		swapData, err := h.svc.BeginSwap(r.Context(), svcQuoteRes)
		res.SwapId = swapData.SwapId.String()
		logctx.Debug(ctx, "BeginSwap", logger.String("swapId", res.SwapId))
		if err != nil {
			restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
			return nil
		}

		signedOrders := []abi.SignedOrder{}

		for i := 0; i < len(swapData.Fragments); i++ {

			// Maker In Amount is Taker's OutAmount!

			// conver In/Out amount to token decimals
			takerInAmount := h.convertToTokenDec(r.Context(), req.InToken, swapData.Fragments[i].OutSize)
			takerOutAmount := h.convertToTokenDec(r.Context(), req.OutToken, swapData.Fragments[i].InSize)

			MakerInAmount := big.NewInt(0)
			MakerInAmount.SetString(takerOutAmount, 10)

			abiOrder := swapData.Orders[i].Signature.AbiFragment
			abiOrder.ExclusivityOverrideBps = big.NewInt(0)

			if len(abiOrder.Outputs) == 0 {
				restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, "abiOrder.Outputs length is 0", logger.Error(err))
				return nil
			}

			// create signed order with amount
			frag := Fragment{
				Signature:      swapData.Orders[i].Signature.Eip712Sig,
				AbiOrder:       abiOrder,
				TakerInAmount:  takerInAmount,
				TakerOutAmount: takerOutAmount,
			}
			res.Fragments = append(res.Fragments, frag)
			// signed order + out amount from the maker's/order side
			signedOrder := abi.SignedOrder{
				OrderWithAmount: abi.OrderWithAmount{
					Order:  abiOrder,
					Amount: MakerInAmount,
				},
				Signature: Signature2Bytes(swapData.Orders[i].Signature.Eip712Sig),
			}
			signedOrders = append(signedOrders, signedOrder)
			// MakerInAmount == takerOutAmount
			logctx.Info(ctx, "append swap fragment", logger.String("swapId", res.SwapId), logger.Int("fragIndex", i), logger.String("TakerInAmount", frag.TakerInAmount), logger.String("MakerInAmount", takerOutAmount))
		}
		// abi encode
		abiCall, err := abi.PackSignedOrders(ctx, signedOrders)
		if err != nil {
			restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
			return nil
		}
		res.AbiCall = fmt.Sprintf("0x%x", abiCall)
		res.Contract = h.reactorAddress
	}

	restutils.WriteJSONResponse(r.Context(), w, http.StatusOK, res)
	return &res
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
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "swapID is empty")
		return nil
	}
	res, err := uuid.Parse(swapId)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "invalid swapID", logger.String("swapId", swapId))
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
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "txHash is empty")
		return
	}

	// execute
	if err := h.svc.SwapStarted(ctx, *swapId, txhash); err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error(), logger.String("swapId", swapId.String()))
		return
	}

	// success
	res := genRes{
		StatusText: "OK",
		Status:     http.StatusOK,
	}
	restutils.WriteJSONResponse(r.Context(), w, http.StatusOK, res, logger.String("swapId started", swapId.String()))

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
		restutils.WriteJSONError(r.Context(), w, http.StatusBadRequest, err.Error(), logger.String("swapId not found", swapId.String()))
		return
	}

	res := genRes{
		StatusText: "OK",
		Status:     http.StatusOK,
	}
	restutils.WriteJSONResponse(r.Context(), w, http.StatusOK, res, logger.String("swapId aborted", swapId.String()))
}
