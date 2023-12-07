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

func TestHandler_CancelOrdersForUser(t *testing.T) {

	t.Run("no user in context - should return `User not found` error", func(t *testing.T) {
		router := mux.NewRouter()

		h, _ := rest.NewHandler(&mocks.MockOrderBookService{}, router)

		req, err := http.NewRequest("DELETE", "/orders", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/orders", h.CancelOrderByOrderId).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "User not found\n", rr.Body.String())
	})

	ctx := mocks.AddUserToCtx(nil)

	tests := []struct {
		name         string
		mockService  *mocks.MockOrderBookService
		expectedCode int
		expectedBody string
	}{
		{
			"no orders found for user",
			&mocks.MockOrderBookService{Error: models.ErrNoOrdersFound},
			http.StatusNotFound,
			"No orders found\n",
		},
		{
			"error cancelling orders for user",
			&mocks.MockOrderBookService{Error: assert.AnError},
			http.StatusInternalServerError,
			"Unable to cancel orders. Try again later\n",
		},
		{
			"successfully cancelled orders for user",
			&mocks.MockOrderBookService{Orders: []models.Order{mocks.Order}},
			http.StatusOK,
			fmt.Sprintf("{\"cancelledOrderIds\":[\"%s\"]}", mocks.Order.Id.String()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			fmt.Println(test.name)

			router := mux.NewRouter()

			h, _ := rest.NewHandler(test.mockService, router)

			req, err := http.NewRequest("DELETE", "/orders", nil)
			if err != nil {
				t.Fatal(err)
			}

			reqWithCtx := req.WithContext(ctx)

			rr := httptest.NewRecorder()
			router.HandleFunc("/orders", h.CancelOrdersForUser).Methods("DELETE")

			router.ServeHTTP(rr, reqWithCtx)

			assert.Equal(t, test.expectedCode, rr.Code)
			assert.Equal(t, test.expectedBody, rr.Body.String())
		})

	}

}
