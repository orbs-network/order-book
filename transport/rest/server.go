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
	// Get all orders for a user
	api.HandleFunc("/orders", PaginationMiddleware(h.GetOrdersForUser)).Methods("GET")
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
	// LH Auction side
	lhApi := h.router.PathPrefix("/lh/v1").Subrouter()
	lhApi.HandleFunc("/begin_auction/{auctionId}", h.beginAuction).Methods("POST")
	lhApi.HandleFunc("/confirm_auction/{auctionId}", h.confirmAuction).Methods("POST")
	lhApi.HandleFunc("/abort_auction/{auctionId}", h.abortAuction).Methods("POST")
	lhApi.HandleFunc("/auction_mined/{auctionId}", h.auctionMined).Methods("POST")

	// unified
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("GET")    // amountOut
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("DELETE") // abort
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("POST")   // confirm
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("PUT")    // mined

	// LISTEN
	logctx.Info(context.TODO(), "starting server", logger.String("port", "8080"))

	if err := http.ListenAndServe(":8080", h.router); err != nil {
		log.Fatalf("error starting http listener: %v", err)
	}
}
