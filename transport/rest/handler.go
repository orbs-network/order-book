package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/shopspring/decimal"
)

// Service represents the methods available on the service to handle the actual request.
type Service interface {
	GetStore() service.OrderBookStore
	ProcessOrder(ctx context.Context, input service.ProcessOrderInput) (models.Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID, isClientOId bool) (cancelledOrderId *uuid.UUID, err error)
	GetBestPriceFor(ctx context.Context, symbol models.Symbol, side models.Side) (decimal.Decimal, error)
	GetOrderById(ctx context.Context, orderId uuid.UUID) (*models.Order, error)
	GetOrderByClientOId(ctx context.Context, clientOId uuid.UUID) (*models.Order, error)
	GetMarketDepth(ctx context.Context, symbol models.Symbol, depth int) (models.MarketDepth, error)
	ConfirmAuction(ctx context.Context, auctionId uuid.UUID) (service.ConfirmAuctionRes, error)
	RevertAuction(ctx context.Context, auctionId uuid.UUID) error
	AuctionMined(ctx context.Context, auctionId uuid.UUID) error
	GetSymbols(ctx context.Context) ([]models.Symbol, error)
	GetOrdersForUser(ctx context.Context, userId uuid.UUID) (orders []models.Order, totalOrders int, err error)
	CancelOrdersForUser(ctx context.Context, publicKey string) error
	GetAmountOut(ctx context.Context, auctionID uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error)
}

type Handler struct {
	svc    Service
	Router *mux.Router
}

func NewHandler(svc Service, r *mux.Router) (*Handler, error) {
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

func (h *Handler) Init() {

	/////////////////////////////////////////////////////////////////////
	// Market maker side
	api := h.Router.PathPrefix("/api/v1").Subrouter()

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
	lhApi := h.Router.PathPrefix("/lh/v1").Subrouter()
	lhApi.HandleFunc("/begin_auction/{auctionId}", h.beginAuction).Methods("POST")
	lhApi.HandleFunc("/confirm_auction/{auctionId}", h.confirmAuction).Methods("GET")
	lhApi.HandleFunc("/abort_auction/{auctionId}", h.abortAuction).Methods("POST")
	lhApi.HandleFunc("/auction_mined/{auctionId}", h.auctionMined).Methods("POST")

	// unified
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("GET")    // amountOut
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("DELETE") // abort
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("POST")   // confirm
	// lhApi.HandleFunc("/auction/{auctionId}", h.amountOut).Methods("PUT")    // mined

	// LISTEN
	//logctx.Info(context.TODO(), "starting server", logger.String("port", port))
	// if err := h.server.ListenAndServe(":8080", h.router); err != nil {
	// 	log.Fatalf("error starting http listener: %v", err)
	// }
}
