package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orbs-network/order-book/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateUserMiddleware(t *testing.T) {

	validApiKey := "Bearer some-api-key"

	tests := []struct {
		name           string
		apiKey         string
		getUserFunc    GetUserByApiKeyFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "should return error if no api key header",
			apiKey:         "",
			getUserFunc:    func(ctx context.Context, apiKey string) (*models.User, error) { return nil, nil },
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid API key (ensure the format is 'Bearer YOUR-API-KEY')\n",
		},
		{
			name:           "should return error if no 'Bearer' keyword",
			apiKey:         "some-api-key",
			getUserFunc:    func(ctx context.Context, apiKey string) (*models.User, error) { return nil, nil },
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid API key (ensure the format is 'Bearer YOUR-API-KEY')\n",
		},
		{
			name:           "should return error if user not found",
			apiKey:         validApiKey,
			getUserFunc:    func(ctx context.Context, apiKey string) (*models.User, error) { return nil, models.ErrUserNotFound },
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized\n",
		},
		{
			name:           "should return error if unexpected error getting user by public key",
			apiKey:         validApiKey,
			getUserFunc:    func(ctx context.Context, apiKey string) (*models.User, error) { return nil, assert.AnError },
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal server error\n",
		},
		{
			name:           "should set user in context when user found",
			apiKey:         validApiKey,
			getUserFunc:    func(ctx context.Context, apiKey string) (*models.User, error) { return &models.User{}, nil },
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

			if tc.apiKey != "" {
				req.Header.Set("X-API-KEY", tc.apiKey)
			}

			rr := httptest.NewRecorder()

			middleware := ValidateUserMiddleware(tc.getUserFunc)
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
