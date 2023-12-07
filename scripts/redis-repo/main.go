package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/redisrepo"
	"github.com/orbs-network/order-book/data/storeuser"
	"github.com/orbs-network/order-book/models"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "createAuctions",
				Usage: "Create auctions",
				Action: func(c *cli.Context) error {
					createAuctions()
					return nil
				},
			},
			{
				Name:  "createOrders",
				Usage: "Create orders",
				Action: func(c *cli.Context) error {
					createOrders()
					return nil
				},
			},
			{
				Name:  "removeOrders",
				Usage: "Remove orders",
				Action: func(c *cli.Context) error {
					removeOrders()
					return nil
				},
			},
			{
				Name:  "createUser",
				Usage: "Create a new user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "apiKey",
						Usage:    "user API key",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					apiKey := c.String("apiKey")

					if apiKey == "" {
						log.Fatalf("apiKey is empty")
					}

					createUser(apiKey)
					return nil
				},
			},
			{
				Name:  "getUserByApiKey",
				Usage: "Get user by their API key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "apiKey",
						Usage:    "user API key",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					apiKey := c.String("apiKey")

					if apiKey == "" {
						log.Fatalf("apiKey is empty")
					}

					getUserByApiKey(apiKey)
					return nil
				},
			},
			{
				Name:  "getUserById",
				Usage: "Get user by their ID",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "user ID",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					userId := c.String("id")

					if userId == "" {
						log.Fatalf("userId is empty")
					}

					getUserById(userId)
					return nil
				},
			},
			{
				Name:  "updateUser",
				Usage: "Updates a user's API key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "newApiKey",
						Usage:    "The desired new user API key",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					newApiKey := c.String("newApiKey")

					if newApiKey == "" {
						log.Fatalf("newApiKey is empty")
					}

					updateUser(newApiKey)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var redisAddress = "localhost:6379"
var redisPassword = "secret"
var redisDb = 0

var rdb = redis.NewClient(&redis.Options{
	Addr:     redisAddress,
	Password: redisPassword,
	DB:       redisDb,
})

var repository, _ = redisrepo.NewRedisRepository(rdb)

var ctx = context.Background()
var orderId = uuid.New()

var userId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var publicKey = "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var size, _ = decimal.NewFromString("10000324.123456789")
var symbol, _ = models.StrToSymbol("USDC-ETH")
var price = decimal.NewFromFloat(10.0)

var user = models.User{
	Id:     userId,
	Type:   models.MARKET_MAKER,
	PubKey: publicKey,
}

func createAuctions() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	fakeAuctionOne := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(200.0),
	}

	fakeAuctionTwo := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(300.0),
	}

	fakeAuctionThree := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(400.0),
	}

	fillOrders := []models.OrderFrag{fakeAuctionOne, fakeAuctionTwo, fakeAuctionThree}

	auctionID := uuid.New()

	ctx := context.Background()

	err = repository.StoreAuction(ctx, auctionID, fillOrders)
	if err != nil {
		log.Fatalf("error storing auction: %v", err)
	}

	auction, err := repository.GetAuction(ctx, auctionID)
	if err != nil {
		log.Fatalf("error getting auction: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("auction: %v", auction)
	log.Print(auction[0].OrderId)
	log.Print(auction[0].Size)
}

func createOrders() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	var order = models.Order{
		Id:        orderId,
		ClientOId: clientOId,
		UserId:    userId,
		Price:     price,
		Symbol:    symbol,
		Size:      size,
		Status:    models.STATUS_OPEN,
		Side:      models.BUY,
		Timestamp: time.Now().UTC(),
	}

	err = repository.StoreOrder(ctx, order)
	if err != nil {
		log.Fatalf("error storing order: %v", err)
	}
}

func removeOrders() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	err = repository.CancelOrdersForUser(ctx, userId)
	if err != nil {
		log.Fatalf("error removing orders: %v", err)
	}
}

func createUser(apiKey string) {
	user.ApiKey = apiKey
	user, err := repository.CreateUser(ctx, user)
	if err != nil {
		log.Fatalf("error storing user: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("userId: %v", user.Id)
	log.Printf("userType: %v", user.Type)
	log.Printf("userPubKey: %v", user.PubKey)
	log.Printf("userApiKey: %v", user.ApiKey)
	log.Print("--------------------------")
}

func getUserByApiKey(apiKey string) {
	retrievedUser, err := repository.GetUserByApiKey(ctx, apiKey)
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("userId: %v", retrievedUser.Id)
	log.Printf("userType: %v", retrievedUser.Type)
	log.Printf("userPubKey: %v", retrievedUser.PubKey)
	log.Print("--------------------------")
}

func getUserById(userIdFlag string) {
	userId := uuid.MustParse(userIdFlag)
	retrievedUser, err := repository.GetUserById(ctx, userId)
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("userId: %v", retrievedUser.Id)
	log.Printf("userType: %v", retrievedUser.Type)
	log.Printf("userPubKey: %v", retrievedUser.PubKey)
	log.Print("--------------------------")
}

func updateUser(newApiKey string) {

	err := repository.UpdateUser(ctx, storeuser.UpdateUserInput{
		UserId: userId,
		PubKey: publicKey,
		ApiKey: newApiKey,
	})

	if err != nil {
		log.Fatalf("error updating user: %v", err)
	}

	log.Printf("user %q updated", userId)
}
