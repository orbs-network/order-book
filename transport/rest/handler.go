package rest

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/middleware"
)

type Handler struct {
	svc      service.OrderBookService
	pairMngr *models.PairMngr
	Router   *mux.Router
}

func NewHandler(svc service.OrderBookService, r *mux.Router) (*Handler, error) {
	if svc == nil {
		return nil, fmt.Errorf("svc cannot be nil")
	}

	if r == nil {
		return nil, fmt.Errorf("router cannot be nil")
	}

	return &Handler{
		svc:      svc,
		Router:   r,
		pairMngr: models.NewPairMngr(),
	}, nil
}

func (h *Handler) Init() {

	/////////////////////////////////////////////////////////////////////
	// Market maker side
	mmApi := h.Router.PathPrefix("/api/v1").Subrouter()

	middlewareValidUser := middleware.ValidateUserMiddleware(h.svc)

	mmApi.Use(middlewareValidUser)

	// ------- CREATE -------
	// Place a new order
	mmApi.HandleFunc("/order", h.ProcessOrder).Methods("POST")

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

	/////////////////////////////////////////////////////////////////////
	// LH Auction side
	lhApi := h.Router.PathPrefix("/lh/v1").Subrouter()
	lhApi.HandleFunc("/begin_auction/{auctionId}", h.beginAuction).Methods("POST")
	lhApi.HandleFunc("/confirm_auction/{auctionId}", h.confirmAuction).Methods("GET")
	lhApi.HandleFunc("/abort_auction/{auctionId}", h.abortAuction).Methods("POST")
	lhApi.HandleFunc("/auction_mined/{auctionId}", h.auctionMined).Methods("POST")

	/////////////////////////////////////////////////////////////////////
	// LH TAKER side -  to replace auction
	takerApi := h.Router.PathPrefix("/taker/v1").Subrouter()
	// returns potential amountOut
	takerApi.HandleFunc("/quote", h.quote).Methods("GET")
	// returns fresh amountOut
	// locks orders
	// returns swapID to be used by abort and txsend
	takerApi.HandleFunc("/swap", h.swap).Methods("GET")
	// release locked orders of start to be used by other match
	// called when
	// lh doesnt want to use swap amountOut
	takerApi.HandleFunc("/abort/{swapId}", h.abortAuction).Methods("POST")
	// all swaps are confirmed on-chain TXHASH
	takerApi.HandleFunc("/txsend/{swapId}", h.auctionMined).Methods("POST")
}
