package rest_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

const ETH_USD = "ETH-USD"

var httpServer *rest.HTTPServer

func runAuctionServer(t *testing.T) {
	repository := mocks.CreateAuctionMock()

	service, err := service.New(repository)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	router := mux.NewRouter()
	handler, err := rest.NewHandler(service, router)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}
	handler.Init()

	httpServer = rest.NewHTTPServer(":8080", handler.Router)
	httpServer.StartServer()
}

type BeginAuctionTest struct {
	name      string
	amountIn  string
	amountOut string
	side      string
	symbol    string
}

func TestHandlers_BeginAuction(t *testing.T) {
	runAuctionServer(t)

	entireA := strconv.Itoa((1) + (2) + (3))
	entireAskB := strconv.Itoa((1000 * 1) + (1001 * 2) + (1002 * 3))
	entireBidB := strconv.Itoa((900 * 1) + (800 * 2) + (700 * 3))

	tests := []BeginAuctionTest{
		{
			name:      "Happy Path BUY - should return 1 ETH for 1000 USD",
			amountIn:  "1000",
			amountOut: "1",
			symbol:    ETH_USD,
			side:      "BUY",
		},
		{
			name:      "Happy Path BUY 2 Orders - should return 2 ETH for 2001 USD",
			amountIn:  "2001",
			amountOut: "2",
			symbol:    ETH_USD,
			side:      "BUY",
		},
		{
			name:      "Partial fill BUY - should return 0.501 ETH for 501 USD",
			amountIn:  "501",
			amountOut: "0.501",
			symbol:    ETH_USD,
			side:      "BUY",
		},
		{
			name:      fmt.Sprintf("EntireBook BUY - should return %s ETH for %s USD", entireA, entireAskB),
			amountIn:  entireAskB,
			amountOut: entireA,
			symbol:    ETH_USD,
			side:      "BUY",
		},
		{
			name:      "Happy Path SELL - should return 900 USD for 1 ETH",
			amountIn:  "1",
			amountOut: "900",
			symbol:    ETH_USD,
			side:      "SELL",
		},
		{
			name:      "Happy Path SELL 2 orders - should return 900+800 USD for 2 ETH",
			amountIn:  "2",
			amountOut: "1700",
			symbol:    ETH_USD,
			side:      "SELL",
		},
		{
			name:      "Partial fill SELL - should return 451 USD for 0.451 ETH",
			amountIn:  "0.5",
			amountOut: "450",
			symbol:    ETH_USD,
			side:      "SELL",
		},
		{
			name:      fmt.Sprintf("EntireBook SELL - should return %s USD for %s ETH", entireBidB, entireA),
			amountIn:  entireA,
			amountOut: entireBidB,
			symbol:    ETH_USD,
			side:      "SELL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := rest.BeginAuctionReq{
				AmountIn: test.amountIn,
				Symbol:   test.symbol,
				Side:     test.side,
			}
			auctionId := uuid.New().String()

			expectedRes := rest.BeginAuctionRes{
				AuctionId: auctionId,
				AmountOut: test.amountOut,
			}

			url := fmt.Sprintf("http://localhost:8080/lh/v1/begin_auction/%s", auctionId)
			jsonData, err := json.Marshal(req)
			assert.NoError(t, err)

			response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Decode the response body into the struct
			var actualRes rest.BeginAuctionRes
			err = json.NewDecoder(response.Body).Decode(&actualRes)
			assert.NoError(t, err)
			assert.Equal(t, expectedRes, actualRes)
		})
	}
	// liquidity insufficient
	t.Run("BUY- should error insuficinet liquidity try to buy with too many B token", func(t *testing.T) {
		insuficientAskB := strconv.Itoa((1000 * 1) + (1001 * 2) + (1002 * 3) + 1)

		req := rest.BeginAuctionReq{
			AmountIn: insuficientAskB,
			Symbol:   ETH_USD,
			Side:     "BUY",
		}
		auctionId := uuid.New().String()
		url := fmt.Sprintf("http://localhost:8080/lh/v1/begin_auction/%s", auctionId)
		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)

		res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		// BAD reQUEST
		assert.Equal(t, res.StatusCode, 400)

		// Read the response body line by line
		defer res.Body.Close()
		reader := bufio.NewReader(res.Body)
		line, err := reader.ReadString('\n')
		assert.NoError(t, err)
		expected := "not enough liquidity in book to satisfy amountIn\n"
		assert.Equal(t, line, expected)
	})
	t.Run("SELL- should error insuficinet liquidity when sell with too many A token", func(t *testing.T) {

		insuficientBidA := strconv.Itoa((1 + 2 + 3) + 1)
		req := rest.BeginAuctionReq{
			AmountIn: insuficientBidA,
			Symbol:   ETH_USD,
			Side:     "SELL",
		}
		auctionId := uuid.New().String()
		url := fmt.Sprintf("http://localhost:8080/lh/v1/begin_auction/%s", auctionId)
		jsonData, err := json.Marshal(req)
		assert.NoError(t, err)

		res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		// BAD reQUEST
		assert.Equal(t, res.StatusCode, 400)

		// Read the response body line by line
		defer res.Body.Close()
		reader := bufio.NewReader(res.Body)
		line, err := reader.ReadString('\n')
		assert.NoError(t, err)
		expected := "not enough liquidity in book to satisfy amountIn\n"
		assert.Equal(t, line, expected)
	})
	// stop server
	httpServer.StopServer(context.Background())
}

func TestHandlers_ConfirmAuction(t *testing.T) {
	runAuctionServer(t)
	// revert auction mock

	t.Run("Happy Path", func(t *testing.T) {
		auctionId := uuid.New().String()

		url := fmt.Sprintf("http://localhost:8080/lh/v1/confirm_auction/%s", auctionId)
		response, err := http.Get(url)
		assert.NoError(t, err)

		// Decode the response body into the struct
		var actualRes rest.ConfirmAuctionRes
		err = json.NewDecoder(response.Body).Decode(&actualRes)
		assert.NoError(t, err)
		assert.Equal(t, len(actualRes.Fragments), 3)
		assert.Equal(t, actualRes.Fragments[0].AmountOut, "1")
		assert.Equal(t, actualRes.Fragments[1].AmountOut, "2")
		assert.Equal(t, actualRes.Fragments[2].AmountOut, "1.5")

		//assert.Equal(t, expectedRes, actualRes)
	})
	// stop server
	httpServer.StopServer(context.Background())
}

func TestHandlers_AbortAuction(t *testing.T) {
	t.Skip("TODO refactor")
	runAuctionServer(t)

	t.Run("Happy Path", func(t *testing.T) {
		auctionId := uuid.New().String()
		url := fmt.Sprintf("http://localhost:8080/lh/v1/abort_auction/%s", auctionId)
		response, err := http.Post(url, "application/json", nil)
		assert.NoError(t, err)
		assert.Equal(t, response.ContentLength, int64(0))
		assert.Equal(t, response.StatusCode, 200)
	})

	// stop server
	httpServer.StopServer(context.Background())
}

func TestHandler_AuctionMined(t *testing.T) {
	runAuctionServer(t)

	t.Run("Auction Mined - but never began or confirmed", func(t *testing.T) {
		auctionId := uuid.New().String()

		url := fmt.Sprintf("http://localhost:8080/lh/v1/auction_mined/%s", auctionId)
		res, err := http.Post(url, "application/json", nil)
		fmt.Println(res)
		fmt.Println(err)
		// Read the response body line by line
		// assert.NoError(t, err)
		// defer res.Body.Close()
		// reader := bufio.NewReader(res.Body)
		// line, err := reader.ReadString('\n')
		// assert.NoError(t, err)
		// expected := "orders in the auction can not fill any longer\n"
		// assert.Equal(t, line, expected)
	})
}
