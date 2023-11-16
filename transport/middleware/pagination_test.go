package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/utils"
	"github.com/stretchr/testify/assert"
)

func TestRest_Pagination(t *testing.T) {
	t.Run("NewPaginationResponse", func(t *testing.T) {
		data := "test"
		page := 1
		pageSize := 10
		total := 100
		totalPages := 10
		ctx := mocks.AddPaginationToCtx(page, pageSize)

		pagination := NewPaginationResponse[string](ctx, data, total)

		assert.Equal(t, data, pagination.Data)
		assert.Equal(t, page, pagination.Page)
		assert.Equal(t, pageSize, pagination.PageSize)
		assert.Equal(t, total, pagination.Total)
		assert.Equal(t, totalPages, pagination.TotalPages)
	})

	t.Run("PaginationMiddleware", func(t *testing.T) {
		testCases := []struct {
			name             string
			url              string
			expectedPage     int
			expectedPageSize int
		}{
			{"no params", "/test", 1, DEFAULT_PAGE_SIZE},
			{"valid params", "/test?page=2&pageSize=20", 2, 20},
			{"invalid page", "/test?page=abc&pageSize=20", 1, 20},
			{"invalid pageSize", "/test?page=2&pageSize=abc", 2, DEFAULT_PAGE_SIZE},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				api := mux.NewRouter()
				var capturedRequest *http.Request

				// Mock handler to capture the request
				testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedRequest = r
				})

				api.HandleFunc("/test", PaginationMiddleware(testHandler)).Methods("GET")

				req, err := http.NewRequest("GET", tt.url, nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				api.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusOK, rr.Code)

				paginator := utils.GetPaginationCtx(capturedRequest.Context())
				assert.Equal(t, tt.expectedPage, paginator.Page)
				assert.Equal(t, tt.expectedPageSize, paginator.PageSize)
			})
		}

	})

}
