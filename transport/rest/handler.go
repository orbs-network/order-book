package rest

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/middleware"
)

type Handler struct {
	svc    service.OrderBookService
	Router *mux.Router
}

func NewHandler(svc service.OrderBookService, r *mux.Router) (*Handler, error) {
	if svc == nil {
		return nil, fmt.Errorf("svc cannot be nil")
	}

	if r == nil {
		return nil, fmt.Errorf("router cannot be nil")
	}

	return &Handler{
		svc:    svc,
		Router: r,
	}, nil
}

func (h *Handler) Init(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	h.initMMRoutes(getUserByApiKey)
	h.initLHRoutes()
}

// Market Maker specific routes
func (h *Handler) initMMRoutes(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	mmApi := h.Router.PathPrefix("/api/v1").Subrouter()

	// Middleware to validate user by API key
	middlewareValidUser := middleware.ValidateUserMiddleware(getUserByApiKey)
	mmApi.Use(middlewareValidUser)

	// ------- CREATE -------
	// Place a new order
	mmApi.HandleFunc("/order", h.CreateOrder).Methods("POST")

	// ------- READ -------
	// Get an order by client order ID
	mmApi.HandleFunc("/order/client-order/{clientOId}", h.GetOrderByClientOId).Methods("GET")
	// Get the best price for a symbol and side
	mmApi.HandleFunc("/order/{side}/{symbol}", h.GetBestPriceFor).Methods("GET")
	// Get an order by ID
	mmApi.HandleFunc("/order/{orderId}", h.GetOrderById).Methods("GET")
	// Get all orders for a user
	mmApi.HandleFunc("/orders", middleware.PaginationMiddleware(h.GetOrdersForUser)).Methods("GET")
	// Get all symbols
	mmApi.HandleFunc("/symbols", h.GetSymbols).Methods("GET")
	// Get market depth
	mmApi.HandleFunc("/orderbook/{symbol}", h.GetMarketDepth).Methods("GET")

	// ------- DELETE -------
	// Cancel an existing order by client order ID
	mmApi.HandleFunc("/order/client-order/{clientOId}", h.CancelOrderByClientOId).Methods("DELETE")
	// Cancel an existing order by order ID
	mmApi.HandleFunc("/order/{orderId}", h.CancelOrderByOrderId).Methods("DELETE")
	// Cancel all orders for a user
	mmApi.HandleFunc("/orders", h.CancelOrdersForUser).Methods("DELETE")
}

// Liquidity Hub specific routes
func (h *Handler) initLHRoutes() {
	lhApi := h.Router.PathPrefix("/lh/v1").Subrouter()

	lhApi.HandleFunc("/begin_auction/{auctionId}", h.beginAuction).Methods("POST")
	lhApi.HandleFunc("/confirm_auction/{auctionId}", h.confirmAuction).Methods("GET")
	lhApi.HandleFunc("/abort_auction/{auctionId}", h.abortAuction).Methods("POST")
	lhApi.HandleFunc("/auction_mined/{auctionId}", h.auctionMined).Methods("POST")
}
