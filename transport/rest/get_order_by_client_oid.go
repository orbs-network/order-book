package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) GetOrderByClientOId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientOIdStr := vars["clientOId"]

	clientOId, err := uuid.Parse(clientOIdStr)
	if err != nil {
		http.Error(w, "invalid clientOId", http.StatusBadRequest)
		return
	}

	logctx.Info(r.Context(), "user trying to get order by clientOId", logger.String("clientOId", clientOId.String()))
	order, err := h.svc.GetOrderByClientOId(r.Context(), clientOId)

	if err != nil {
		http.Error(w, "Internal error. Try again later", http.StatusInternalServerError)
		return
	}

	if order == nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(order)

	if err != nil {
		logctx.Error(r.Context(), "error marshaling order", logger.Error(err))
		http.Error(w, "Error getting order by client ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		logctx.Error(r.Context(), "failed to write response", logger.Error(err))
		http.Error(w, "Error getting order by client ID", http.StatusInternalServerError)
	}

}
