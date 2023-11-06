package rest_test

import (
	"bufio"
	"bytes"
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

type ErrorResponse struct {
	Message string `json:"message"`
}

func createServer(t *testing.T) bool {
	auctionRepo := mocks.CreateAuctionMock()
	service, err := service.New(auctionRepo)
	//service, err := service.New(repository)
	if err != nil {
		t.Fatalf("error creating service: %v", err)
		return false
	}

	router := mux.NewRouter()

	handler, err := rest.NewHandler(service, router)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}

	go handler.Listen()

	return true

}

type BeginAuctionTest struct {
	name      string
	amountIn  string
	amountOut string
	side      string
	symbol    string
}

func TestHandler_beginAuction(t *testing.T) {

	res := createServer(t)
	assert.True(t, res)

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
	t.Run("begin_auction BUY- liquidity insuficinet try to buy with too many B token", func(t *testing.T) {
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
	t.Run("begin_auction BUY- liquidity insuficinet try to buy with too many B token", func(t *testing.T) {

		insuficientBidB := strconv.Itoa((900 * 1) + (800 * 2) + (700 * 3) + 1)
		req := rest.BeginAuctionReq{
			AmountIn: insuficientBidB,
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
}
