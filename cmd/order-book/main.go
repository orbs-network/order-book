package main

import (
	"log"
	"os"
	"strconv"

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
	redisAddress, found := os.LookupEnv("REDIS_ADDRESS")
	if !found {
		panic("REDIS_ADDRESS not set")
	}

	redisPassword, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		panic("REDIS_PASSWORD not set")
	}

	redisDb, found := os.LookupEnv("REDIS_DB")
	if !found {
		panic("REDIS_DB not set")
	}

	redisDbInt, err := strconv.Atoi(redisDb)
	if err != nil {
		panic("REDIS_DB not a number")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       redisDbInt,
	})

	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	service, err := service.New(repository)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	router := mux.NewRouter()

	handler, err := rest.NewHandler(service, router)
	if err != nil {
		log.Fatalf("error creating handler: %v", err)
	}

	handler.Listen()
}
