package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func writeJSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, resp interface{}, logFields ...logger.Field) {
	// Buffer the response body
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(resp); err != nil {
		logctx.Error(ctx, "failed to encode JSON response", logFields...)
		http.Error(w, "Error processing request. Try again later", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
