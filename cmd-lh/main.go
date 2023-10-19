package main

import (
	"context"
	"log"

	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

func main() {
	setup()
}

func setup() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "secret",
		DB:       0,
	})

	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	s, err := service.New(repository)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	ctx := context.Background()
	symbol, _ := models.StrToSymbol("ETH-USD")
	// id1 := uuid.New()
	// id2 := uuid.New()
	// id3 := uuid.New()
	// ASK
	// s.AddOrder(ctx, id3, decimal.NewFromInt(1002), symbol, decimal.NewFromInt(3), models.SELL)
	// s.AddOrder(ctx, id2, decimal.NewFromInt(1001), symbol, decimal.NewFromInt(2), models.SELL)
	// s.AddOrder(ctx, id1, decimal.NewFromInt(1000), symbol, decimal.NewFromInt(1), models.SELL)

	// BIDS
	// s.AddOrder(ctx, id1, decimal.NewFromInt(900), symbol, decimal.NewFromInt(1), models.BUY)
	// s.AddOrder(ctx, id2, decimal.NewFromInt(800), symbol, decimal.NewFromInt(2), models.BUY)
	// s.AddOrder(ctx, id3, decimal.NewFromInt(700), symbol, decimal.NewFromInt(3), models.BUY)

	// Amount for ask
	amountIn := 1000 + 1001 + 1001 + 1002 + 1002 + 501
	amountOut, err := s.GetAmountOut(ctx, symbol, models.BUY, decimal.NewFromInt(int64(amountIn)))
	if err != nil {
		log.Fatalf("error GetAmountOut: %v", err)
	}
	// should get 5.5
	log.Printf("Amount out is: %s", amountOut.String())

	// Amount for ask bids
	amountIn = 1 + 1
	amountOut, err = s.GetAmountOut(ctx, symbol, models.SELL, decimal.NewFromInt(int64(amountIn)))
	if err != nil {
		log.Fatalf("error GetAmountOut: %v", err)
	}
	// should get 900 + 800 =1700
	log.Printf("Amount out is: %s", amountOut.String())

	// Not working...
	// s.CancelOrder(ctx, id1)
	// s.CancelOrder(ctx, id2)
	// s.CancelOrder(ctx, id3)

}
