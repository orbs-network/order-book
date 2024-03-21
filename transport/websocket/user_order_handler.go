package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/middleware"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

// WebSocketOrderHandler returns a handler that upgrades the connection to WebSocket and subscribes to order updates for a particular user
// The user is authenticated using the API key in the request
func WebSocketOrderHandler(orderSvc service.OrderBookService, getUserByApiKey middleware.GetUserByApiKeyFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract API key from query parameters
		apiKey, err := middleware.BearerToken(r, "X-API-KEY")
		if err != nil {
			logctx.Warn(r.Context(), "incorrect API key format", logger.Error(err))
			http.Error(w, "Invalid API key (ensure the format is 'Bearer YOUR-API-KEY')", http.StatusBadRequest)
			return
		}

		// Authenticate user
		user, err := getUserByApiKey(r.Context(), apiKey)
		if err != nil {
			logctx.Warn(r.Context(), "user not found by api key")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Upgrade to WebSocket
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logctx.Error(r.Context(), "error upgrading to websocket", logger.Error(err))
			http.Error(w, "Error subscribing to orders", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Subscribe to that user's order updates
		messageChan, err := orderSvc.SubscribeUserOrders(r.Context(), user.Id)
		if err != nil {
			logctx.Error(r.Context(), "error subscribing to user orders", logger.Error(err))
			http.Error(w, "Error subscribing to orders", http.StatusInternalServerError)
			return
		}

		// Read messages from the channel and send to WebSocket
		for msg := range messageChan {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logctx.Error(r.Context(), "error writing to websocket", logger.Error(err))
				break
			}

		}
	}
}
