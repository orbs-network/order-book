package restutils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func WriteJSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, resp interface{}, logFields ...logger.Field) {
	// Buffer the response body
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(resp); err != nil {
		logFields = append(logFields, logger.Error(err))
		logctx.Error(ctx, "failed to encode JSON response", logFields...)
		WriteJSONError(w, http.StatusInternalServerError, "Error processing request. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(buf.Bytes()); err != nil {
		logFields = append(logFields, logger.Error(err))
		logctx.Error(ctx, "failed to write response", logFields...)
		WriteJSONError(w, http.StatusInternalServerError, "Error processing request. Try again later")
	}
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResponse := ErrorResponse{
		Status: status,
		Msg:    message,
	}

	json.NewEncoder(w).Encode(errResponse)
}
