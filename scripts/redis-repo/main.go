package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/data/redisrepo"
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
				Name:  "storeUserByPubKey",
				Usage: "Store user by public key",
				Action: func(c *cli.Context) error {
					storeUserByPublicKey()
					return nil
				},
			},
			{
				Name:  "getUserByPubKey",
				Usage: "Get user by public key",
				Action: func(c *cli.Context) error {
					getUserByPublicKey()
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

func storeUserByPublicKey() {
	err := repository.StoreUserByPublicKey(ctx, user)
	if err != nil {
		log.Fatalf("error storing user: %v", err)
	}
}

func getUserByPublicKey() {
	retrievedUser, err := repository.GetUserByPublicKey(ctx, publicKey)
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}
	log.Print("--------------------------")
	log.Printf("userId: %v", retrievedUser.Id)
	log.Printf("userType: %v", retrievedUser.Type)
	log.Printf("userPubKey: %v", retrievedUser.PubKey)
	log.Print("--------------------------")
}
