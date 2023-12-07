package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/rest"
)

func main() {
	setup()
}

func setup() {
	redisAddress, found := os.LookupEnv("REDIS_URL")
	if !found {
		panic("REDIS_URL not set")
	}

	opt, err := redis.ParseURL(redisAddress)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis url: %v", err))
	}

	log.Printf("Redis address: %s", opt.Addr)

	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8080"
	}

	rdb := redis.NewClient(opt)

	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	// TODO: add CLI flag to easily switch between blockchains
	ethClient := &service.EthereumClient{}

	service, err := service.New(repository, ethClient)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	router := mux.NewRouter()
	handler, err := rest.NewHandler(service, router)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}
	handler.Init()

	server := rest.NewHTTPServer(":"+port, handler.Router)
	server.StartServer()
	// blocking
	<-server.StopChannel
	//handler.Listen(port)
}
