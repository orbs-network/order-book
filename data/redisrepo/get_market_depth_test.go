package redisrepo

// func TestRedisRepository_GetMarketDepth(t *testing.T) {
// 	ctx := context.Background()
// 	symbol := "BTC-ETH"
// 	depth := 10

// 	db, mock := redismock.NewClientMock()

// 	repo := &redisRepository{
// 		client: db,
// 	}

// 	// Generate some test data
// 	for i := 0; i < depth; i++ {
// 		order := models.Order{
// 			ID:             uuid.New(),
// 			Price:          decimal.NewFromFloat(100.0),
// 			AvailableSize:  decimal.NewFromFloat(1.0),
// 			RemainingSize:  decimal.NewFromFloat(1.0),
// 			Side:           models.SideBuy,
// 			Status:         models.StatusOpen,
// 			CreationTime:   time.Now(),
// 			LastUpdateTime: time.Now(),
// 		}
// 		err := repo.CreateOrder(ctx, order)
// 		if err != nil {
// 			t.Fatalf("Failed to create test order: %v", err)
// 		}
// 	}

// 	// Call the GetMarketDepth function
// 	marketDepth, err := repo.GetMarketDepth(ctx, symbol, depth)
// 	if err != nil {
// 		t.Fatalf("Failed to get market depth: %v", err)
// 	}

// 	// Verify the market depth
// 	if len(marketDepth.Asks) != depth {
// 		t.Errorf("Unexpected number of asks. Expected: %d, Got: %d", depth, len(marketDepth.Asks))
// 	}
// 	if len(marketDepth.Bids) != depth {
// 		t.Errorf("Unexpected number of bids. Expected: %d, Got: %d", depth, len(marketDepth.Bids))
// 	}
// }
