package rest

import (
	"context"
	"encoding/json"
	"fmt"
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

func (h *Handler) ToTokenBigInt(ctx context.Context, tokenName string, amount decimal.Decimal) *big.Int {
	if token := h.supportedTokens.ByName(tokenName); token != nil {
		dcmls := decimal.NewFromInt(10)
		dcmls = dcmls.Pow(decimal.NewFromInt(int64(token.Decimals)))
		mul := amount.Mul(dcmls)

		// get rid of decimal point values (round down)
		// remove after point 1.00123 decimals
		mul = mul.Floor()
		// convert
		bgint := big.NewInt(0)
		bgint.SetString(mul.String(), 10)

		return bgint
	}
	logctx.Error(ctx, "Token is not found in supported tokens: "+tokenName)
	return nil
}

func (h *Handler) convertFromTokenDec(ctx context.Context, tokenName, amountStr string) (decimal.Decimal, error) {
	if token := h.supportedTokens.ByName(tokenName); token != nil {
		dcmls := decimal.NewFromInt(10)
		dcmls = dcmls.Pow(decimal.NewFromInt(int64(token.Decimals)))
		amount, err := decimal.NewFromString(amountStr)
		if err != nil {
			logctx.Error(ctx, "error converting amountStr: "+amountStr)
			return decimal.Zero, err
		}
		res := amount.Div(dcmls)
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
func (h *Handler) nameFromAddress(address string) (string, error) {
	token := h.supportedTokens.ByAddress(address)
	if token == nil {
		return "", models.ErrTokenNotsupported
	}
	return token.Name, nil
}

func (h *Handler) resolveQuoteTokenNames(req *QuoteReq) error {
	// has address but no name
	if req.InToken == "" {
		InName, err := h.nameFromAddress(req.InTokenAddress)
		if err != nil {
			return err
		}
		if len(InName) > 0 {
			req.InToken = InName
		}
	}
	if req.OutToken == "" {
		OutName, err := h.nameFromAddress(req.OutTokenAddress)
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
		logctx.Warn(ctx, "handleQuote Failed to decode json", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error())
		return nil
	}

	logctx.Debug(ctx, "QuoteReq", logger.String("InToken", req.InToken), logger.String("InAmount", req.InAmount), logger.String("OutToken", req.OutToken), logger.String("MinOutAmount", req.MinOutAmount))

	// ensure token names if only addresses were sent
	err = h.resolveQuoteTokenNames(&req)
	if err != nil {
		logctx.Warn(ctx, "handleQuote Failed to resolveQuoteTokenNames", logger.Error(err))
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error(), logger.String("InTokenAddress", req.InTokenAddress), logger.String("OutTokenAddress", req.OutTokenAddress))
		return nil
	}

	inAmount, err := h.convertFromTokenDec(ctx, req.InToken, req.InAmount)

	if err != nil {
		logctx.Warn(ctx, "handleQuote Failed to convertFromTokenDec", logger.Error(err))
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
	// taker's in token to maker's side
	makerSide := pair.GetMakerSide(req.InToken)

	// resolve makerInAddress to verify balance on-chain
	makerInAdrs := req.OutTokenAddress
	if makerInAdrs == "" {
		makerInAdrs = h.supportedTokens.ByName(req.OutToken).Address
	}

	// ALWAYS reverese decimals tp meet the makers order's side
	svcQuoteRes, err := h.svc.GetQuote(r.Context(), pair.Symbol(), makerSide, inAmount, minOutAmount, makerInAdrs)
	if err != nil {
		if err == models.ErrMinOutAmount {
			restutils.WriteJSONError(ctx, w, http.StatusBadRequest, err.Error())
		} else if err == models.ErrInsufficientBalance {
			restutils.WriteJSONError(ctx, w, http.StatusConflict, err.Error())
		} else {
			restutils.WriteJSONError(ctx, w, http.StatusInternalServerError, err.Error())
		}
		return nil
	}

	convOutAmount := h.ToTokenBigInt(r.Context(), req.OutToken, svcQuoteRes.Size)
	if convOutAmount == nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "convOutAmount return empty string")
		return nil
	}
	// convert res
	res := QuoteRes{
		OutAmount: convOutAmount.String(),
		OutToken:  req.OutToken,
		InAmount:  req.InAmount,
		InToken:   req.InToken,
		//SwapId:    "",
		Fragments: []Fragment{},
	}

	logctx.Debug(ctx, "QuoteRes OK", logger.String("OutAmount", res.OutAmount))

	if isSwap {
		// lock liquidity
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
			// convert to sol's big int and floor (reduce precision here)
			takerInAmount := h.ToTokenBigInt(r.Context(), req.InToken, swapData.Fragments[i].InSize)
			takerOutAmount := h.ToTokenBigInt(r.Context(), req.OutToken, swapData.Fragments[i].OutSize)

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
				TakerInAmount:  takerInAmount.String(),
				TakerOutAmount: takerOutAmount.String(),
			}
			res.Fragments = append(res.Fragments, frag)
			// signed order + out amount from the maker's/order side
			signedOrder := abi.SignedOrder{
				OrderWithAmount: abi.OrderWithAmount{
					Order:  abiOrder,
					Amount: takerInAmount, // is what the taker requested to swap for this frag
				},
				Signature: Signature2Bytes(swapData.Orders[i].Signature.Eip712Sig),
			}
			signedOrders = append(signedOrders, signedOrder)
			// MakerInAmount == takerOutAmount
			logctx.Debug(ctx, "append swap fragment", logger.String("swapId", res.SwapId), logger.Int("fragIndex", i), logger.String("TakerInAmount", frag.TakerInAmount), logger.String("takerOutAmount", takerOutAmount.String()))
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
	logctx.Info(r.Context(), "abortSwap called from POST")

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
