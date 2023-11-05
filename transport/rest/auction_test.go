package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/rest"
	"github.com/stretchr/testify/assert"
)

func createServer(t *testing.T) bool {
	//address := "localhost:6379"
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     address,
	// 	Password: "secret",
	// 	DB:       10, // for test
	// })
	// if rdb == nil {
	// 	fmt.Println("redis is not running in ", address)
	// 	return false
	// }
	// repository, err := redisrepo.NewRedisRepository(rdb)
	// if err != nil {
	// 	t.Fatalf("error creating repository: %v", err)
	// }

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
func TestHandler_beginAuction(t *testing.T) {

	res := createServer(t)
	assert.True(t, res)
	req := rest.BeginAuctionReq{
		AmountIn: "1000",
		Symbol:   "ETH-USD",
		Side:     "BUY",
	}
	auctionId := uuid.New().String()

	expectedRes := rest.BeginAuctionRes{
		AuctionId: auctionId,
		AmountOut: "1",
	}

	// expected_res := AmountOutResponse{}

	// Send a GET request to the URL

	url := fmt.Sprintf("http://localhost:8080/lh/v1/begin_auction/%s", auctionId)
	jsonData, err := json.Marshal(req)
	assert.NoError(t, err)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Decode the response body into the struct
	actualRes := rest.BeginAuctionRes{}
	err = json.NewDecoder(response.Body).Decode(&actualRes)
	assert.NoError(t, err)
	assert.Equal(t, expectedRes, actualRes)

}
