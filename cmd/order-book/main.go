package main

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/data"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/rest"
)

func main() {
	setup()
}

func setup() {
	repository, err := data.NewMemoryRespository()
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
