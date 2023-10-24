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

	// ------- CREATE -------
	// Place a new order
	api.HandleFunc("/order", h.ProcessOrder).Methods("POST")

	// ------- READ -------
	// Get an order by client order ID
	api.HandleFunc("/order/client-order/{clientOId}", h.GetOrderByClientOId).Methods("GET")
	// Get the best price for a symbol and side
	api.HandleFunc("/order/{side}/{symbol}", h.GetBestPriceFor).Methods("GET")
	// Get an order by ID
	api.HandleFunc("/order/{orderId}", h.GetOrderById).Methods("GET")
	// Get all symbols
	api.HandleFunc("/symbols", h.GetSymbols).Methods("GET")
	// Get market depth
	api.HandleFunc("/orderbook/{symbol}", h.GetMarketDepth).Methods("GET")

	// ------- DELETE -------
	// Cancel an existing order by client order ID
	api.HandleFunc("/order/client-order/{clientOId}", h.CancelOrderByClientOId).Methods("DELETE")
	// Cancel an existing order by order ID
	api.HandleFunc("/order/{orderId}", h.CancelOrderByOrderId).Methods("DELETE")

	/////////////////////////////////////////////////////////////////////
	// LH side
	lhApi := h.router.PathPrefix("/lh/v1").Subrouter()
	lhApi.HandleFunc("/quote", h.amountOut).Methods("POST")
	lhApi.HandleFunc("/approve_orders", h.approveOrders).Methods("GET")

	// LISTEN
	logctx.Info(context.TODO(), "starting server", logger.String("port", "8080"))

	if err := http.ListenAndServe(":8080", h.router); err != nil {
		log.Fatalf("error starting http listener: %v", err)
	}
}
