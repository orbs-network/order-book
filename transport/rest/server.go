package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (h *Handler) Listen() {

	/////////////////////////////////////////////////////////////////////
	// Market maker side
	api := h.router.PathPrefix("/api/v1").Subrouter()
	// Create a new order
	api.HandleFunc("/order", h.ProcessOrder).Methods("POST")
	// Cancel an existing order
	api.HandleFunc("/order/{orderId}", h.CancelOrder).Methods("DELETE")
	// Get the best price for a symbol and side
	api.HandleFunc("/order/{side}/{symbol}", h.GetBestPriceFor).Methods("GET")
	// Get an order by id
	api.HandleFunc("/order/{orderId}", h.GetOrderById).Methods("GET")
	// Get market depth
	api.HandleFunc("/orderbook/{symbol}", h.GetMarketDepth).Methods("GET")

	/////////////////////////////////////////////////////////////////////
	// LH side
	lhApi := h.router.PathPrefix("/lh/v1").Subrouter()
	lhApi.HandleFunc("/quote", h.amountOut).Methods("POST")
	lhApi.HandleFunc("/confirm_auction", h.confirmAuction).Methods("GET")

	// LISTEN
	logctx.Info(context.TODO(), "starting server", logger.String("port", "8080"))

	if err := http.ListenAndServe(":8080", h.router); err != nil {
		log.Fatalf("error starting http listener: %v", err)
	}
}
