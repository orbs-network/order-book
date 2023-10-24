package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetSymbols(t *testing.T) {
	router := mux.NewRouter()

	symbol := models.Symbol("TRX-USDC")
	symbols := []models.Symbol{}
	symbols = append(symbols, symbol)

	h, _ := rest.NewHandler(&mocks.MockOrderBookService{Symbols: symbols}, router)

	req, err := http.NewRequest("GET", "/symbols", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.HandleFunc("/symbols", h.GetSymbols).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), fmt.Sprintf(`"symbol":"%s"`, symbol))
}
