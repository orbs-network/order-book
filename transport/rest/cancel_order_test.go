package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CancelOrderByOrderId(t *testing.T) {
	var orderId = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	tests := []struct {
		name         string
		mockService  *mocks.MockOrderBookService
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			"invalid orderId",
			&mocks.MockOrderBookService{},
			"/order/invalid",
			http.StatusBadRequest,
			"invalid order ID\n",
		},
		{
			"no user in context",
			&mocks.MockOrderBookService{Error: models.ErrNoUserInContext},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusUnauthorized,
			"User not found\n",
		},
		{
			"order not found",
			&mocks.MockOrderBookService{Error: models.ErrOrderNotFound},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusNotFound,
			"Order not found\n",
		},
		{
			"trying to cancel order that is not open",
			&mocks.MockOrderBookService{Error: models.ErrOrderNotOpen},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusNotFound,
			"Order not found\n",
		},
		{
			"unexpected error from service",
			&mocks.MockOrderBookService{Error: assert.AnError},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusInternalServerError,
			"Error cancelling order. Try again later\n",
		},
		{
			"no cancelledOrderId returned from service",
			&mocks.MockOrderBookService{},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusInternalServerError,
			"Error cancelling order. Try again later\n",
		},
		{
			"successful cancel",
			&mocks.MockOrderBookService{Order: &models.Order{Id: orderId}},
			fmt.Sprintf("/order/%s", orderId.String()),
			http.StatusOK,
			fmt.Sprintf("{\"orderId\":\"%s\"}", orderId.String()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			fmt.Println(test.name)

			router := mux.NewRouter()

			h, _ := rest.NewHandler(test.mockService, router)

			req, err := http.NewRequest("DELETE", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.HandleFunc("/order/{orderId}", h.CancelOrderByOrderId).Methods("DELETE")

			router.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedCode, rr.Code)
			assert.Equal(t, test.expectedBody, rr.Body.String())
		})

	}
}

func TestHandler_CancelOrderByClientOId(t *testing.T) {
	var orderId = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000009")

	tests := []struct {
		name         string
		mockService  *mocks.MockOrderBookService
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			"invalid clientOId",
			&mocks.MockOrderBookService{},
			"/order/client-order/invalid",
			http.StatusBadRequest,
			"invalid clientOId\n",
		},
		{
			"no user in context",
			&mocks.MockOrderBookService{Error: models.ErrNoUserInContext},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusUnauthorized,
			"User not found\n",
		},
		{
			"order not found",
			&mocks.MockOrderBookService{Error: models.ErrOrderNotFound},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusNotFound,
			"Order not found\n",
		},
		{
			"trying to cancel order that is not open",
			&mocks.MockOrderBookService{Error: models.ErrOrderNotOpen},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusNotFound,
			"Order not found\n",
		},
		{
			"unexpected error from service",
			&mocks.MockOrderBookService{Error: assert.AnError},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusInternalServerError,
			"Error cancelling order. Try again later\n",
		},
		{
			"no cancelledOrderId returned from service",
			&mocks.MockOrderBookService{},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusInternalServerError,
			"Error cancelling order. Try again later\n",
		},
		{
			"successful cancel",
			&mocks.MockOrderBookService{Order: &models.Order{Id: orderId}},
			fmt.Sprintf("/order/client-order/%s", clientOId.String()),
			http.StatusOK,
			fmt.Sprintf("{\"orderId\":\"%s\"}", orderId.String()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			fmt.Println(test.name)

			router := mux.NewRouter()

			h, _ := rest.NewHandler(test.mockService, router)

			req, err := http.NewRequest("DELETE", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router.HandleFunc("/order/client-order/{clientOId}", h.CancelOrderByClientOId).Methods("DELETE")

			router.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedCode, rr.Code, fmt.Sprintf("expected code %d, got %d", test.expectedCode, rr.Code))
			assert.Equal(t, test.expectedBody, rr.Body.String(), fmt.Sprintf("expected body %s, got %s", test.expectedBody, rr.Body.String()))
		})

	}
}
