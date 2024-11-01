package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	// Handle SIGINT and SIGTERM signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown handler
	go func() {
		sig := <-signalChan
		log.Printf("Tracker: received SIGTERM signal: %s. Initiating shutdown...\n", sig)
		cancel()
	}()

	log.Printf("Swaps tracker running with ticker duration: %s\n", tickerDuration)

	for {
		select {
		case <-ticker.C:
			// updating maker on-chain wallets per token
			err := evmClient.UpdateMakerBalances(ctx)
			if err != nil {
				log.Printf("error checking makers token balance: %v\n", err)
			}

			// check pending transactions
			err = evmClient.CheckPendingTxs(ctx)
			if err != nil {
				log.Printf("error checking pending txs: %v", err)
			}

		case <-ctx.Done():
			log.Println("Shutting down swaps tracker...")
			return
		}
	}
}
