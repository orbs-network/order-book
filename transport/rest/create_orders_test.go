package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

func createOrdersBody(t *testing.T, args rest.CreateOrdersRequest) *bytes.Buffer {
	json, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(json)
}

func TestHandler_CreateOrders(t *testing.T) {
	orderReqOne := rest.CreateOrderRequest{
		Price:         "100.0",
		Size:          "10",
		Symbol:        "MATIC-USDC",
		Side:          "sell",
		ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
		Eip712Sig:     "mock-sig",
		Eip712Msg:     mocks.MsgData,
	}

	orderReqTwo := rest.CreateOrderRequest{
		Price:         "102.0",
		Size:          "20",
		Symbol:        "MATIC-USDC",
		Side:          "sell",
		ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
		Eip712Sig:     "mock-sig",
		Eip712Msg:     mocks.MsgData,
	}

	ctx := mocks.AddUserToCtx(nil)

	t.Run("mismatching symbols are not allowed", func(t *testing.T) {
		router := mux.NewRouter()

		h, _ := rest.NewHandler(&mocks.MockOrderBookService{}, router)

		mismatchingSymbolOrder := rest.CreateOrderRequest{
			Price:         "102.0",
			Size:          "20",
			Symbol:        "BTC-USDC",
			Side:          "sell",
			ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			Eip712Sig:     "mock-sig",
			Eip712Msg:     mocks.MsgData,
		}

		body := createOrdersBody(t, rest.CreateOrdersRequest{
			Symbol: "MATIC-USDC",
			Orders: []rest.CreateOrderRequest{orderReqOne, mismatchingSymbolOrder},
		})

		req, err := http.NewRequest("POST", "/orders", body)
		if err != nil {
			t.Fatal(err)
		}

		reqWithCtx := req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router.HandleFunc("/orders", h.CreateOrders).Methods("POST")

		router.ServeHTTP(rr, reqWithCtx)

		assert.Equal(t, 400, rr.Code)
		assert.Equal(t, "{\"symbol\":\"MATIC-USDC\",\"created\":[],\"status\":400,\"msg\":\"Symbol in order \\\"BTC-USDC\\\" does not match symbol in request \\\"MATIC-USDC\\\"\"}\n", rr.Body.String())
	})

	t.Run("valid request", func(t *testing.T) {
		router := mux.NewRouter()

		h, _ := rest.NewHandler(&mocks.MockOrderBookService{
			Order: &models.Order{},
		}, router)

		body := createOrdersBody(t, rest.CreateOrdersRequest{
			Symbol: "MATIC-USDC",
			Orders: []rest.CreateOrderRequest{orderReqOne, orderReqTwo},
		})

		req, err := http.NewRequest("POST", "/orders", body)
		if err != nil {
			t.Fatal(err)
		}

		reqWithCtx := req.WithContext(ctx)

		rr := httptest.NewRecorder()
		router.HandleFunc("/orders", h.CreateOrders).Methods("POST")

		router.ServeHTTP(rr, reqWithCtx)

		assert.Equal(t, 201, rr.Code)
	})
}
