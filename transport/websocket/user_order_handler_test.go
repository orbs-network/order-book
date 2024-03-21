package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
)

func TestWebSocketOrderHandler(t *testing.T) {

	t.Run("Test successful websocket lifecycle", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebSocketOrderHandler(&mocks.MockOrderBookService{}, mockGetUserByApiKey)(w, r)
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[len("http"):]

		dialer := websocket.Dialer{}
		headers := http.Header{}
		headers.Set("X-API-KEY", "Bearer mock-api-key")

		conn, _, err := dialer.Dial(wsURL, headers)
		assert.NoError(t, err)
		defer conn.Close()

		assert.NotNil(t, conn, "The WebSocket connection should be established")

		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		assert.NoError(t, err, "Should be able to close the WebSocket connection")

		conn.Close()
	})

	t.Run("Test invalid API key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebSocketOrderHandler(&mocks.MockOrderBookService{}, mockGetUserByApiKeyError)(w, r)
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[len("http"):]

		dialer := websocket.Dialer{}
		headers := http.Header{}
		headers.Set("X-API-KEY", "Bearer invalid-api-key")

		conn, _, err := dialer.Dial(wsURL, headers)
		assert.Error(t, err)
		assert.Nil(t, conn, "The WebSocket connection should not be established")
	})

	t.Run("Test error upgrading to WebSocket", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebSocketOrderHandler(&mocks.MockOrderBookService{}, mockGetUserByApiKey)(w, r)
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[len("http"):]

		dialer := websocket.Dialer{}
		headers := http.Header{}
		headers.Set("X-API-KEY", "Bearer mock-api-key")

		conn, _, err := dialer.Dial(wsURL, headers)
		assert.NoError(t, err)
		defer conn.Close()

		conn.Close()

		_, _, err = conn.ReadMessage()
		assert.Error(t, err, "Expect an error after the connection is closed")
	})

	t.Run("Test error subscribing to user orders", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WebSocketOrderHandler(&mocks.MockOrderBookService{Error: assert.AnError}, mockGetUserByApiKey)(w, r)
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[len("http"):]

		dialer := websocket.Dialer{}
		headers := http.Header{}
		headers.Set("X-API-KEY", "Bearer mock-api-key")

		conn, _, err := dialer.Dial(wsURL, headers)
		assert.NoError(t, err)
		defer conn.Close()

		_, _, err = conn.ReadMessage()
		assert.Error(t, err, "Expect an error after the connection is closed")
	})
}

var mockGetUserByApiKey = func(ctx context.Context, apiKey string) (*models.User, error) {
	return &models.User{
		Id:     uuid.MustParse("00000000-0000-0000-0000-000000000007"),
		ApiKey: "mock-api-key",
	}, nil
}

var mockGetUserByApiKeyError = func(ctx context.Context, apiKey string) (*models.User, error) {
	return nil, assert.AnError
}
