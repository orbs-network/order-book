package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserMiddleware(t *testing.T) {

	tests := []struct {
		name           string
		publicKey      string
		mockService    mocks.MockOrderBookService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "should return error if no public key header",
			publicKey:      "",
			mockService:    mocks.MockOrderBookService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing public key\n",
		},
		{
			name:           "should return error if user not found",
			publicKey:      "some-public-key",
			mockService:    mocks.MockOrderBookService{Error: models.ErrUserNotFound},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized\n",
		},
		{
			name:           "should return error if unexpected error getting user by public key",
			publicKey:      "some-public-key",
			mockService:    mocks.MockOrderBookService{Error: assert.AnError},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error\n",
		},
		{
			name:      "should set user in context when user found",
			publicKey: "some-public-key",
			mockService: mocks.MockOrderBookService{
				User: &mocks.User},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.publicKey != "" {
				req.Header.Set("X-Public-Key", tc.publicKey)
			}

			rr := httptest.NewRecorder()

			middleware := ValidateUserMiddleware(&tc.mockService)
			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// This is the next handler in the chain, it would normally process the request
				// if the middleware passes it through
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			assert.Equal(t, tc.expectedBody, rr.Body.String())
		})
	}
}
