package rest

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) fakeFill(w http.ResponseWriter, r *http.Request) {
	res := h.handleQuote(w, r, true)
	if res == nil {
		logctx.Warn(r.Context(), "FakeFill Failed")
		return
	}
	logctx.Debug(r.Context(), "FakeFill", logger.String("swapId", res.SwapId), logger.String("InAmount", res.InAmount), logger.String("OutAmount", res.OutAmount))
	err := h.svc.FillSwap(r.Context(), uuid.MustParse(res.SwapId))
	if err != nil {
		restutils.WriteJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}
	// response is written in FillSwap
}
