package rest_test

import (
	"bytes"
	"encoding/json"
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

func createBody(t *testing.T, args rest.CreateOrderRequest) *bytes.Buffer {
	orderReqJSON, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewBuffer(orderReqJSON)
}

func TestHandler_CreateOrder(t *testing.T) {

	orderReq := rest.CreateOrderRequest{
		Price:         "100.0",
		Size:          "10",
		Symbol:        "BTC-USD",
		Side:          "sell",
		ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
	}

	orderReqJSON, _ := json.Marshal(orderReq)

	orderSucessRes := rest.CreateOrderResponse{
		OrderId: mocks.Order.Id.String(),
	}

	orderSucessResJSON, _ := json.Marshal(orderSucessRes)

	t.Run("no user in context - should return `User not found` error", func(t *testing.T) {
		router := mux.NewRouter()

		h, _ := rest.NewHandler(&mocks.MockOrderBookService{}, router)

		req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(orderReqJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order", h.CreateOrder).Methods("POST")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "User not found\n", rr.Body.String())
	})

	ctx := mocks.AddUserToCtx(nil)

	tests := []struct {
		name         string
		mockService  *mocks.MockOrderBookService
		body         *bytes.Buffer
		expectedCode int
		expectedBody any
	}{
		// ----- Required fields validation tests -----
		{
			"no price in request body - should return `missing required field 'price'` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Size:          "10",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"missing required field 'price'\n",
		},
		{
			"no size in request body - should return `missing required field 'size'` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"missing required field 'size'\n",
		},
		{
			"no symbol in request body - should return `missing required field 'symbol'` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "10",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"missing required field 'symbol'\n",
		},
		{
			"no side in request body - should return `missing required field 'side'` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "10",
				Symbol:        "BTC-USD",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"missing required field 'side'\n",
		},
		{
			"no client order id in request body - should return `missing required field 'clientOrderId'` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:  "100.0",
				Size:   "10",
				Symbol: "BTC-USD",
				Side:   "sell",
			}),
			http.StatusBadRequest,
			"missing required field 'clientOrderId'\n",
		},
		// ----- Parse fields tests -----
		{
			"invalid price format - should return `price is not a valid number format` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0.0",
				Size:          "10",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'price' is not a valid number format\n",
		},
		{
			"negative price - should return `price must be positive` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "-100.0",
				Size:          "10",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'price' must be positive\n",
		},
		{
			"invalid size format - should return `size is not a valid number format` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "dsfdsfsdf",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'size' is not a valid number format\n",
		},
		{
			"negative size - should return `size must be positive` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "-10",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'size' must be positive\n",
		},
		{
			"invalid symbol - should return `symbol is not valid` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "10",
				Symbol:        "BTC-SOME-INVALID-SYMBOL",
				Side:          "sell",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'symbol' is not valid\n",
		},
		{
			"invalid side - should return `side is not valid` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "10",
				Symbol:        "BTC-USD",
				Side:          "some-invalid-side",
				ClientOrderId: "a677273e-12de-4acc-a4f8-de7fb5b86e37",
			}),
			http.StatusBadRequest,
			"'side' is not valid\n",
		},
		{
			"invalid client order id - should return `clientOrderId is not valid` error",
			&mocks.MockOrderBookService{},
			createBody(t, rest.CreateOrderRequest{
				Price:         "100.0",
				Size:          "10",
				Symbol:        "BTC-USD",
				Side:          "sell",
				ClientOrderId: "1",
			}),
			http.StatusBadRequest,
			"'clientOrderId' is not valid\n",
		},
		// ----- Create order tests -----
		{
			"create order success - should return `order created`",
			&mocks.MockOrderBookService{Order: &mocks.Order},
			createBody(t, orderReq),
			http.StatusCreated,
			string(orderSucessResJSON),
		},
		{
			"clashing order id - should return `Clashing order ID. Please retry` error",
			&mocks.MockOrderBookService{Order: &models.Order{}, Error: models.ErrClashingOrderId},
			createBody(t, orderReq),
			http.StatusConflict,
			"Clashing order ID. Please retry\n",
		},
		{
			"clashing client order id - should return `Order with clientOrderId %q already exists. You must first cancel this order` error",
			&mocks.MockOrderBookService{Order: &models.Order{}, Error: models.ErrClashingClientOrderId},
			createBody(t, orderReq),
			http.StatusConflict,
			fmt.Sprintf("Order with clientOrderId %q already exists. You must first cancel this order\n", orderReq.ClientOrderId),
		},
		{
			"unexpected error - should return `Error creating order. Try again later` error",
			&mocks.MockOrderBookService{Order: &models.Order{}, Error: assert.AnError},
			createBody(t, orderReq),
			http.StatusInternalServerError,
			"Error creating order. Try again later\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			fmt.Println(test.name)

			router := mux.NewRouter()

			h, _ := rest.NewHandler(test.mockService, router)

			req, err := http.NewRequest("POST", "/order", test.body)
			if err != nil {
				t.Fatal(err)
			}

			reqWithCtx := req.WithContext(ctx)

			rr := httptest.NewRecorder()
			router.HandleFunc("/order", h.CreateOrder).Methods("POST")

			router.ServeHTTP(rr, reqWithCtx)

			assert.Equal(t, test.expectedCode, rr.Code)
			assert.Equal(t, test.expectedBody, rr.Body.String())
		})

	}

}
