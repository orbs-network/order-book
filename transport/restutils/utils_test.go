package restutils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/stretchr/testify/assert"
)

func TestWriteJSONResponse(t *testing.T) {
	type CreateOrdersResponse struct {
		Data          []*models.Order `json:"data"`
		Status        string          `json:"status"`
		FailureReason string          `json:"failureReason"`
	}

	w := httptest.NewRecorder()

	WriteJSONResponse(context.Background(), w, http.StatusBadRequest, CreateOrdersResponse{
		Data:          []*models.Order{},
		Status:        "400",
		FailureReason: "Bad order",
	}, logger.String("test", "test"))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp CreateOrdersResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, CreateOrdersResponse{
		Data:          []*models.Order{},
		Status:        "400",
		FailureReason: "Bad order",
	}, resp)
}

func TestWriteJSONResponse_EncodeError(t *testing.T) {
	type UnsupportedType struct {
		Ch chan int
	}

	w := httptest.NewRecorder()

	WriteJSONResponse(context.Background(), w, http.StatusOK, UnsupportedType{}, logger.String("test", "test"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "{\"status\":500,\"msg\":\"Error processing request. Try again later\"}\n", w.Body.String())
}
