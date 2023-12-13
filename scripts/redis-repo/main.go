// Various scripts for testing the redis repository

package main

import (
	"context"
	"fmt"
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
				Name:  "createSwaps",
				Usage: "Create swaps",
				Action: func(c *cli.Context) error {
					createSwaps()
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
				Name:  "markOrderAsFilled",
				Usage: "Mark order as filled",
				Action: func(c *cli.Context) error {
					markOrderAsFilled()
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
				Name:  "getOrdersByIds",
				Usage: "Get multiple orders by their IDs",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "ids",
						Usage:    "order IDs",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					ids := c.StringSlice("ids")
					fmt.Printf("ids: %#v\n", ids)

					if ids == nil {
						log.Fatalf("ids is empty")
					}

					getOrdersByIds(ids)
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
var orderId = uuid.MustParse("9bfc6d29-07e0-4bf7-9189-bc03bdadb1ae")

var userId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var publicKey = "0x6a04ab98d9e4774ad806e302dddeb63bea16b5cb5f223ee77478e861bb583eb336b6fbcb60b5b3d4f1551ac45e5ffc4936466e7d98f6c7c0ec736539f74691a6"
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var size, _ = decimal.NewFromString("1000")
var symbol, _ = models.StrToSymbol("USDC-ETH")
var price = decimal.NewFromFloat(10.0)

var user = models.User{
	Id:     userId,
	Type:   models.MARKET_MAKER,
	PubKey: publicKey,
}

func createSwaps() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	fakeSwapOne := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(200.0),
	}

	fakeSwapTwo := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(300.0),
	}

	fakeSwapThree := models.OrderFrag{
		OrderId: uuid.New(),
		Size:    decimal.NewFromFloat(400.0),
	}

	fillOrders := []models.OrderFrag{fakeSwapOne, fakeSwapTwo, fakeSwapThree}

	swapId := uuid.New()

	ctx := context.Background()

	err = repository.StoreSwap(ctx, swapId, fillOrders)
	if err != nil {
		log.Fatalf("error storing swap: %v", err)
	}

	swap, err := repository.GetSwap(ctx, swapId)
	if err != nil {
		log.Fatalf("error getting swap: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("swap: %v", swap)
	log.Print(swap[0].OrderId)
	log.Print(swap[0].Size)
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
		Side:      models.BUY,
		Timestamp: time.Now().UTC(),
	}

	err = repository.StoreOpenOrder(ctx, order)
	if err != nil {
		log.Fatalf("error storing order: %v", err)
	}
}

func markOrderAsFilled() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	order, err := repository.FindOrderById(ctx, orderId, false)
	if err != nil {
		log.Fatalf("error getting order: %v", err)
	}

	order.SizeFilled = order.Size

	err = repository.StoreFilledOrders(ctx, []models.Order{*order})
	if err != nil {
		log.Fatalf("error marking order as filled: %v", err)
	}
}

func removeOrders() {
	repository, err := redisrepo.NewRedisRepository(rdb)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	cancelledOrderIds, err := repository.CancelOrdersForUser(ctx, userId)
	if err != nil {
		log.Fatalf("error removing orders: %v", err)
	}

	log.Print("--------------------------")
	log.Printf("cancelledOrderIds: %v", cancelledOrderIds)
	log.Print("--------------------------")
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

func getOrdersByIds(orderIds []string) {
	ids := make([]uuid.UUID, len(orderIds))
	for i, id := range orderIds {
		ids[i] = uuid.MustParse(id)
	}

	orders, err := repository.FindOrdersByIds(ctx, ids)
	if err != nil {
		log.Fatalf("error getting users: %v", err)
	}
	for _, order := range orders {
		log.Print("--------------------------")
		log.Printf("orderId: %v", order.Id)
		log.Printf("orderClientOId: %v", order.ClientOId)
		log.Printf("orderUserId: %v", order.UserId)
		log.Printf("orderPrice: %v", order.Price)
		log.Printf("orderSymbol: %v", order.Symbol)
		log.Printf("orderSize: %v", order.Size)
		log.Printf("orderSide: %v", order.Side)
		log.Printf("orderTimestamp: %v", order.Timestamp)
		log.Print("--------------------------")
	}
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
