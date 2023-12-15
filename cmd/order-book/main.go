package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"github.com/orbs-network/order-book/data/evmrepo"
	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/serviceuser"
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

	rpcUrl, found := os.LookupEnv("RPC_URL")
	if !found {
		panic("RPC_URL not set")
	}

	rdb := redis.NewClient(opt)

	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("error creating eth client: %v", err)
	}
	defer ethClient.Close()

	evmRepo, err := evmrepo.NewEvmRepository(ethClient)
	if err != nil {
		log.Fatalf("error creating evm repository: %v", err)
	}

	// TODO: add CLI flag to easily switch between blockchains
	evmClient, err := service.NewEvmSvc(repository, evmRepo)
	if err != nil {
		log.Fatalf("error creating evm client: %v", err)
	}

	service, err := service.New(repository, evmClient)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	userSvc, err := serviceuser.New(repository)
	if err != nil {
		log.Fatalf("error creating user service: %v", err)
	}

	router := mux.NewRouter()
	handler, err := rest.NewHandler(service, router)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}

	handler.Init(userSvc.GetUserByApiKey)

	server := rest.NewHTTPServer(":"+port, handler.Router)
	server.StartServer()
	// blocking
	<-server.StopChannel
	//handler.Listen(port)
}
