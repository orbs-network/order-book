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

var orderId = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func TestHandler_CancelOrder(t *testing.T) {

	t.Run("invalid orderId - should return `http.StatusBadRequest`", func(t *testing.T) {
		mockService := &mocks.MockOrderBookService{}
		router := mux.NewRouter()

		h, _ := rest.NewHandler(mockService, router)

		req, err := http.NewRequest("DELETE", "/order/invalid", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "invalid order ID\n", rr.Body.String())

	})

	t.Run("no user in context - should return `http.StatusUnauthorized`", func(t *testing.T) {
		mockService := &mocks.MockOrderBookService{Error: models.ErrNoUserInContext}
		router := mux.NewRouter()

		h, _ := rest.NewHandler(mockService, router)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/order/%s", orderId.String()), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "User not found\n", rr.Body.String())
	})

	t.Run("order not found - should return `http.StatusNotFound`", func(t *testing.T) {
		mockService := &mocks.MockOrderBookService{Error: models.ErrOrderNotFound}
		router := mux.NewRouter()

		h, _ := rest.NewHandler(mockService, router)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/order/%s", orderId.String()), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "Order not found\n", rr.Body.String())
	})

	t.Run("unexpected error from service - should return `http.StatusInternalServerError`", func(t *testing.T) {
		mockService := &mocks.MockOrderBookService{Error: assert.AnError}
		router := mux.NewRouter()

		h, _ := rest.NewHandler(mockService, router)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/order/%s", orderId.String()), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Error cancelling order. Try again later\n", rr.Body.String())
	})

	t.Run("successful cancel - should return `http.StatusOK`", func(t *testing.T) {
		mockService := &mocks.MockOrderBookService{}
		router := mux.NewRouter()

		h, _ := rest.NewHandler(mockService, router)

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/order/%s", orderId.String()), nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "{\"orderId\":\"00000000-0000-0000-0000-000000000001\"}", rr.Body.String())
	})

}
