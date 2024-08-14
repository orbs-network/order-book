package websocket

import (
	"net/http"
	"time"

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
		ctx := r.Context()
		// Extract API key from query parameters
		apiKey, err := middleware.BearerToken(r, "X-API-KEY")
		if err != nil {
			logctx.Warn(ctx, "incorrect API key format", logger.Error(err))
			http.Error(w, "Invalid API key (ensure the format is 'Bearer YOUR-API-KEY')", http.StatusBadRequest)
			return
		}

		// Authenticate user
		user, err := getUserByApiKey(ctx, apiKey)
		if err != nil {
			logctx.Warn(ctx, "user not found by api key")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Upgrade to WebSocket
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logctx.Error(ctx, "error upgrading to websocket", logger.Error(err), logger.String("userId", user.Id.String()))
			http.Error(w, "Error subscribing to orders", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		conn.SetPongHandler(func(appData string) error {
			conn.SetReadDeadline(time.Now().Add(120 * time.Second)) // Extend deadline on pong
			return nil
		})

		messageChan, err := orderSvc.SubscribeUserOrders(ctx, user.Id)
		if err != nil {
			logctx.Error(ctx, "error subscribing to user orders", logger.Error(err), logger.String("userId", user.Id.String()))
			http.Error(w, "Error subscribing to orders", http.StatusInternalServerError)
			return
		}

		// Ensure Redis connection is unsubscribed and closed when the WebSocket disconnects
		defer func() {
			orderSvc.UnsubscribeUserOrders(ctx, user.Id, messageChan)
		}()

		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-messageChan:
				if !ok {
					logctx.Warn(ctx, "message channel closed", logger.String("userId", user.Id.String()))
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					logctx.Warn(ctx, "unable to write to websocket", logger.Error(err), logger.String("userId", user.Id.String()))
					return
				}
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					logctx.Error(ctx, "error sending ping", logger.Error(err), logger.String("userId", user.Id.String()))
					return
				}
			case <-ctx.Done():
				logctx.Info(ctx, "request context cancelled", logger.String("userId", user.Id.String()))
				return
			}
		}
	}
}
