package restutils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

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
		WriteJSONError(ctx, w, http.StatusInternalServerError, "Error processing request. Try again later")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(buf.Bytes()); err != nil {
		logFields = append(logFields, logger.Error(err))
		logctx.Error(ctx, "failed to write response", logFields...)
		WriteJSONError(ctx, w, http.StatusInternalServerError, "Error processing request. Try again later")
	}
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func WriteJSONError(ctx context.Context, w http.ResponseWriter, status int, message string, logFields ...logger.Field) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errResponse := ErrorResponse{
		Status: status,
		Msg:    message,
	}

	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		logFields = append(logFields, logger.Error(err))
		logctx.Error(ctx, "failed to write error response", logFields...)
		return
	}

	// log details about why the request failed or was rejected
	logFields = append(logFields, logger.String("sys_msg", message))
	logFields = append(logFields, logger.String("status", http.StatusText(status)))
	logctx.Warn(ctx, "api request not successful", logFields...)
}

// read os env var with default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
