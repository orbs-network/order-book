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
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/serviceuser"
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
			{
				Name:  "storePendingSwap",
				Usage: "Store a pending swap",
				Action: func(c *cli.Context) error {
					storePendingSwap()
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
var usersvc, _ = serviceuser.New(repository)

var ctx = context.Background()
var orderId = uuid.MustParse("9bfc6d29-07e0-4bf7-9189-bc03bdadb1ae")

var userId = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var publicKey = "0x37edf0398ec38921baa65dfc808d151d1bc979c7c3af9649bbde160b96b2851599b9b13fe138116ffafee0ebd775ecc7bcb9ba911aa488c10db3d4a26b72178e"
var clientOId = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var size, _ = decimal.NewFromString("1000")
var symbol, _ = models.StrToSymbol("MATIC-USDC")
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
		OutSize: decimal.NewFromFloat(200.0),
	}

	fakeSwapTwo := models.OrderFrag{
		OrderId: uuid.New(),
		OutSize: decimal.NewFromFloat(300.0),
	}

	fakeSwapThree := models.OrderFrag{
		OrderId: uuid.New(),
		OutSize: decimal.NewFromFloat(400.0),
	}

	fillOrders := []models.OrderFrag{fakeSwapOne, fakeSwapTwo, fakeSwapThree}

	swapId := uuid.New()

	ctx := context.Background()

	err = repository.StoreSwap(ctx, swapId, fillOrders)
	if err != nil {
		log.Fatalf("error storing swap: %v", err)
	}

	swap, err := repository.GetSwap(ctx, swapId, true)
	if err != nil {
		log.Fatalf("error getting swap: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("swap: %v", swap)
	log.Print(swap.Frags[0].OrderId)
	log.Print(swap.Frags[0].OutSize)
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

	mockBcClient := &mocks.MockBcClient{}

	svc, err := service.New(repository, mockBcClient)
	if err != nil {
		log.Fatalf("error creating service: %v", err)
	}

	// NOTE: This ignores pagination
	orders, _, err := repository.GetOpenOrders(ctx, userId)
	if err != nil {
		log.Fatalf("error getting orders: %v", err)
	}

	var cancelledOrderIds []uuid.UUID

	for _, order := range orders {
		_, err = svc.CancelOrder(ctx, service.CancelOrderInput{
			Id:          order.Id,
			IsClientOId: false,
			UserId:      order.UserId,
		})
		if err != nil {
			log.Fatalf("error cancelling order: %v", err)
		}
		cancelledOrderIds = append(cancelledOrderIds, order.Id)
	}

	log.Print("--------------------------")
	log.Printf("cancelledOrderIds: %v", cancelledOrderIds)
	log.Print("--------------------------")
}

func createUser(apiKey string) {
	user.ApiKey = apiKey
	user, apiKey, err := usersvc.CreateUser(ctx, serviceuser.CreateUserInput{
		PubKey: publicKey,
	})
	if err != nil {
		log.Fatalf("error creating user: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("userId: %v", user.Id)
	log.Printf("userType: %v", user.Type)
	log.Printf("userPubKey: %v", user.PubKey)
	log.Printf("hashedApiKey: %v", user.ApiKey)
	log.Printf("apiKey: %v", apiKey)
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

	orders, err := repository.FindOrdersByIds(ctx, ids, false)
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

func storePendingSwap() {
	p := models.SwapTx{
		SwapId: uuid.New(),
		TxHash: "0x5dcbfe934287c50363e5c82502739aadd4d535a1f7c0ccd7a8088fb4dfd800da",
	}

	for i := 0; i < 1000; i++ {
		err := repository.StoreNewPendingSwap(ctx, p)
		if err != nil {
			log.Fatalf("error storing pending swap: %v", err)
		}
	}
}
