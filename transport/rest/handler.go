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
		reactorAddress:  restutils.GetEnv("SWAP_CONTRACT_ADDRESS", "0x2ee46d8d20020520d5266f3cacc7c41e1aadd4c6"),
	}, nil
}

func (h *Handler) Init(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	h.initMakerRoutes(getUserByApiKey)
	h.initTakerRoutes(getUserByApiKey)
}

// Market Maker specific routes
func (h *Handler) initMakerRoutes(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	mmApi := h.Router.PathPrefix("/api/v1").Subrouter()

	disableUserAuth := restutils.GetEnv("DISABLE_USER_AUTH", "")

	// Middleware to validate user by API key
	middlewareValidUser := middleware.ValidateUserMiddleware(getUserByApiKey)
	if disableUserAuth != "true" {
		mmApi.Use(middlewareValidUser)
	}

	// ------- CREATE -------
	// Place multiple orders
	mmApi.HandleFunc("/orders", h.CreateOrders).Methods("POST")
	// Place a new order
	mmApi.HandleFunc("/order", h.CreateOrder).Methods("POST")

	// ------- READ -------
	// Get an order by client order ID
	mmApi.HandleFunc("/order/client-order/{clientOId}", h.GetOrderByClientOId).Methods("GET")
	// Get an order by ID
	mmApi.HandleFunc("/order/{orderId}", h.GetOrderById).Methods("GET")
	// Get all open orders for a user
	mmApi.HandleFunc("/orders", middleware.PaginationMiddleware(h.GetOpenOrdersForUser)).Methods("GET")
	// Get all filled orders for a user
	mmApi.HandleFunc("/fills", middleware.PaginationMiddleware(h.GetFilledOrdersForUser)).Methods("GET")
	// Get all symbols
	mmApi.HandleFunc("/symbols", h.GetSymbols).Methods("GET")
	// Get market depth
	mmApi.HandleFunc("/orderbook/{symbol}", h.GetMarketDepth).Methods("GET")
	// Get supported tokens
	mmApi.HandleFunc("/supported-tokens", h.GetSupportedTokens).Methods("GET")

	// ------- DELETE -------
	// Cancel an existing order by client order ID
	mmApi.HandleFunc("/order/client-order/{clientOId}", h.CancelOrderByClientOId).Methods("DELETE")
	// Cancel an existing order by order ID
	mmApi.HandleFunc("/order/{orderId}", h.CancelOrderByOrderId).Methods("DELETE")
	// Cancel all orders for a user
	mmApi.HandleFunc("/orders", h.CancelOrdersForUser).Methods("DELETE")
}

// Liquidity Hub specific routes
func (h *Handler) initTakerRoutes(getUserByApiKey middleware.GetUserByApiKeyFunc) {
	/////////////////////////////////////////////////////////////////////
	// TAKER side
	takerApi := h.Router.PathPrefix("/taker/v1").Subrouter()

	// disabled!
	//middlewareValidUser := middleware.ValidateUserMiddleware(getUserByApiKey)
	//takerApi.Use(middlewareValidUser) disable for now

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
