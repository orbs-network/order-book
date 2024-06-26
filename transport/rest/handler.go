package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/featureflags"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/middleware"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/transport/websocket"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type Handler struct {
	svc             service.OrderBookService
	pairMngr        *models.PairMngr
	Router          *mux.Router
	okJson          []byte
	supportedTokens *service.SupportedTokens
	reactorAddress  string
}
type genRes struct {
	StatusText string `json:"statusText"`
	Status     int    `json:"status"`
}

func NewHandler(svc service.OrderBookService, r *mux.Router) (*Handler, error) {
	var supportedTokensPath = os.Getenv("SUPPORTED_TOKENS_JSON_FILE_PATH")

	if svc == nil {
		return nil, fmt.Errorf("svc cannot be nil")
	}

	if r == nil {
		return nil, fmt.Errorf("router cannot be nil")
	}

	// Create an empty JSON object
	okJsonObj := genRes{
		StatusText: "OK",
		Status:     http.StatusOK,
	}

	// Convert the emptyJSON object to JSON format
	okJson, err := json.Marshal(okJsonObj)
	if err != nil {
		return nil, err
	}

	if supportedTokensPath == "" {
		logctx.Warn(context.Background(), "SUPPORTED_TOKENS_JSON_FILE_PATH env var not set, using default")
		supportedTokensPath = "supportedTokens.json"
	}

	// load supported tokens
	st, err := service.NewSupportedTokens(context.Background(), supportedTokensPath)
	if st == nil {
		logctx.Error(context.Background(), "failed to load supported tokens", logger.Error(err))
		return nil, err
	}

	return &Handler{
		svc:             svc,
		Router:          r,
		pairMngr:        models.NewPairMngr(),
		okJson:          okJson,
		supportedTokens: st,
		reactorAddress:  restutils.GetEnv("REACTOR_ADDRESS", "0x4C4B950432189b3283A5111A6963ee318109695c"),
	}, nil
}

func (h *Handler) Init(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	h.initMakerRoutes(getUserByApiKey)
	h.initTakerRoutes(getUserByApiKey)
}

// Market Maker specific routes
func (h *Handler) initMakerRoutes(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	mmApi := h.Router.PathPrefix("/api/v1").Subrouter()

	// Middleware to validate user by API key
	middlewareValidUser := middleware.ValidateUserMiddleware(getUserByApiKey)
	mmApi.Use(middlewareValidUser)

	getApi := mmApi.Methods("GET").Subrouter()
	// Only users with write permissions can access these routes
	createApi := mmApi.Methods("POST").Subrouter()
	createApi.Use(middleware.CheckUserHasPermsMiddleware([]models.UserType{"MARKET_MAKER", "ADMIN"}))
	deleteApi := mmApi.Methods("DELETE").Subrouter()
	deleteApi.Use(middleware.CheckUserHasPermsMiddleware([]models.UserType{"MARKET_MAKER", "ADMIN"}))

	// ------- CREATE -------
	// Place multiple orders
	createApi.HandleFunc("/orders", h.CreateOrders).Methods("POST")
	// Place a new order
	createApi.HandleFunc("/order", h.CreateOrder).Methods("POST")

	// ------- READ -------
	// Get an order by client order ID
	getApi.HandleFunc("/order/client-order/{clientOId}", h.GetOrderByClientOId)
	// Get an order by ID
	getApi.HandleFunc("/order/{orderId}", h.GetOrderById)
	// Get all open orders for a user
	getApi.HandleFunc("/orders", middleware.PaginationMiddleware(h.GetOpenOrders))
	// correct way to Get fills using data in swaps
	getApi.HandleFunc("/fills", h.GetSwapFills)
	// Get all symbols
	getApi.HandleFunc("/symbols", h.GetSymbols)
	// Get market depth
	getApi.HandleFunc("/orderbook/{symbol}", h.GetMarketDepth)
	// Get supported tokens
	getApi.HandleFunc("/supported-tokens", h.GetSupportedTokens)

	// ------- DELETE -------
	// Cancel an existing order by client order ID
	deleteApi.HandleFunc("/order/client-order/{clientOId}", h.CancelOrderByClientOId).Methods("DELETE")
	// Cancel an existing order by order ID
	deleteApi.HandleFunc("/order/{orderId}", h.CancelOrderByOrderId).Methods("DELETE")
	// Cancel all orders for a user
	deleteApi.HandleFunc("/orders", h.CancelOrdersForUser).Methods("DELETE")

	// ------- WEBSOCKET -------
	// Subscribe to order events (websocket)
	getApi.HandleFunc("/ws/orders", websocket.WebSocketOrderHandler(h.svc, getUserByApiKey))
}

// Liquidity Hub specific routes
func (h *Handler) initTakerRoutes(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	/////////////////////////////////////////////////////////////////////
	// TAKER side
	takerApi := h.Router.PathPrefix("/taker/v1").Subrouter()

	// disabled!
	middlewareValidUser := middleware.ValidateUserMiddleware(getUserByApiKey)
	takerApi.Use(middlewareValidUser)

	// IN: InAmount, InToken, OutToken or InTokenAddress, OutTokenAddress
	// OUT: CURRENT potential outAmount
	takerApi.HandleFunc("/quote", h.quote).Methods("POST")
	// IN: InAmount, InToken, OutToken or InTokenAddress, OutTokenAddress
	// OUT: Locked outAmount, SwapID
	takerApi.HandleFunc("/swap", h.swap).Methods("POST")
	// IN: SwapID given in /swap
	// IN: txHash
	// start tracking txhash onchain
	takerApi.HandleFunc("/started/{swapId}/{txHash}", h.swapStarted).Methods("POST")
	// IN: SwapID given in /swap
	// release locked orders of start to be used by other match
	// called when lh doesnt want to use swap outAmount
	takerApi.HandleFunc("/abort/{swapId}", h.abortSwap).Methods("POST")
	// IN: txHash, SwapID given in /swap
	// Notifies order book to start tracking the state of the tx (discuss events or based on txHash)
	// takerApi.HandleFunc("/txsent/{swapId}", h.txSent).Methods("POST")

	// ------- DEBUG -------
	if featureflags.FlagEnableFakeFill != "" {
		takerApi.HandleFunc("/fake-fill", h.fakeFill).Methods("POST")
	}
}
