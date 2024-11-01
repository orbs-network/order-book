// A script to create a new market maker user and print the API key
// The stored API key is SHA256 hashed, so it cannot be recovered
// Usage: go run scripts/create-user/main.go
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/serviceuser"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Creating a new read-write user...")

	ctx := context.Background()

	// Retrieve REDIS_URL from environment variable or prompt for it
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Print("Enter REDIS_URL: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		redisURL = scanner.Text()
	}

	// Create a Redis client
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	if strings.HasPrefix(redisURL, "rediss") {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	rdb := redis.NewClient(opt)

	// Instantiate repository and service
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	usersvc, err := serviceuser.New(repository)
	if err != nil {
		log.Fatalf("Failed to create user service: %v", err)
	}

	// Retrieve wallet address from environment variable or prompt for it
	walletAddress := os.Getenv("WALLET_ADDRESS")
	if walletAddress == "" {
		fmt.Print("Enter wallet address: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		walletAddress = scanner.Text()
	}

	// Create a new user
	user, apiKey, err := usersvc.CreateUser(ctx, serviceuser.CreateUserInput{
		PubKey:   walletAddress,
		UserType: "MARKET_MAKER",
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Print the created user and API key
	fmt.Printf("User Created: %+v\n", user)
	fmt.Println("----------------------------------------")
	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Println("⚠️⚠️⚠️ The API key is SHA256 hashed and cannot be recovered ⚠️⚠️⚠️")
	fmt.Println("----------------------------------------")

	err = writeToFile("api_key.txt", apiKey)
	if err != nil {
		log.Fatalf("Failed to write API key to file: %v", err)
	}
	fmt.Println("API key written to api_key.txt")

}

func writeToFile(filename, text string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}
