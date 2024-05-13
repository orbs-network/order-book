package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/orbs-network/order-book/data/evmrepo"
	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/service"
	"github.com/redis/go-redis/v9"
)

var defaultDuration = 10 * time.Second

func main() {
	log.Printf("Starting pending swaps tracker...")

	redisAddress, found := os.LookupEnv("REDIS_URL")
	if !found {
		redisAddress, found = os.LookupEnv("REDISCLOUD_URL")
		if !found {
			panic("Neither REDIS_URL nor REDISCLOUD_URL is set")
		}
	}

	opt, err := redis.ParseURL(redisAddress)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis url: %v", err))
	}

	log.Printf("Redis address: %s", opt.Addr)

	rpcUrl, found := os.LookupEnv("RPC_URL")
	if !found {
		panic("RPC_URL not set")
	}

	if strings.HasPrefix(redisAddress, "rediss") {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
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

	evmClient, err := service.NewEvmSvc(repository, evmRepo)
	if err != nil {
		log.Fatalf("error creating evm client: %v", err)
	}

	envDurationStr := os.Getenv("TICKER_DURATION")

	tickerDuration, err := time.ParseDuration(envDurationStr)
	if err != nil || envDurationStr == "" {
		fmt.Printf("Invalid or missing TICKER_DURATION. Using default of %s\n", defaultDuration)
		tickerDuration = defaultDuration
	}

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	ctx := context.Background()

	for range ticker.C {
		// updating maker on chain wallets per token
		err := evmClient.UpdateMakerBalance(ctx)
		if err != nil {
			log.Printf("error checking makers token balance: %v", err)
		}

		// check pending transactions
		err = evmClient.CheckPendingTxs(ctx)
		if err != nil {
			log.Printf("error checking pending txs: %v", err)
		}
	}
}
