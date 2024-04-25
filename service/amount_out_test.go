package service

// /////////////////////////////////////////////////////////////////
// static iterator impl
// type iter struct {
// 	orders []models.Order
// 	index  int
// }

// func (i *iter) HasNext() bool {
// 	return i.index < (len(i.orders) - 1)
// }
// func (i *iter) Next(ctx context.Context) *models.Order {
// 	// increment index
// 	i.index = i.index + 1

// 	if i.index >= len(i.orders) {
// 		return nil
// 	}

// 	// get order
// 	return &i.orders[i.index]
// }

// func newOrder(price, size int64) models.Order {
// 	return models.Order{
// 		Price: decimal.NewFromInt(price),
// 		Size:  decimal.NewFromInt(size),
// 	}
// }
// func newAsks() models.OrderIter {
// 	return &iter{
// 		orders: []models.Order{
// 			newOrder(1000, 1),
// 			newOrder(1001, 2),
// 			newOrder(1002, 3),
// 		},
// 		index: -1,
// 	}
// }
// func newBids() models.OrderIter {
// 	return &iter{
// 		orders: []models.Order{
// 			newOrder(900, 1),
// 			newOrder(800, 2),
// 			newOrder(700, 3),
// 		},
// 		index: -1,
// 	}
// }

// /////////////////////////////////////////////////////////////////
// func TestService_getOutAmount(t *testing.T) {
// 	ctx := context.Background()

// 	// order have no onchain data hence -1 as decimals are ignored in impl

// 	aTokenDec := int32(18)
// 	//bTokenDec := int32(6)
// 	t.Run("getOutAmountInAToken- happy path", func(t *testing.T) {
// 		res, err := getOutAmountInAToken(ctx, newAsks(), decimal.NewFromFloat((1000*1)+(1001*2)+(1002*3)))
// 		assert.Equal(t, err, nil)
// 		assert.Equal(t, res.Size.Round(aTokenDec).String(), decimal.NewFromFloat(1+2+3).Round(aTokenDec).String())
// 	})
// 	t.Run("getOutAmountInAToken- Partial fill", func(t *testing.T) {
// 		res, err := getOutAmountInAToken(ctx, newAsks(), decimal.NewFromFloat(501))
// 		assert.Equal(t, err, nil)
// 		assert.Equal(t, res.Size.String(), decimal.NewFromFloat(0.501).String())
// 	})

// 	t.Run("getOutAmountInAToken- liquidity insufficient", func(t *testing.T) {
// 		_, err := getOutAmountInAToken(ctx, newAsks(), decimal.NewFromFloat((1000*1)+(1001*2)+(1002*3)+1))
// 		assert.Equal(t, err, models.ErrInsufficientLiquity)
// 	})

// 	t.Run("getOutAmountInBToken- happy path", func(t *testing.T) {
// 		res, err := getOutAmountInBToken(ctx, newBids(), decimal.NewFromFloat(1+2+3))
// 		assert.Equal(t, err, nil)
// 		assert.Equal(t, res.Size.String(), decimal.NewFromFloat((900*1)+(800*2)+(700*3)).String())
// 	})

// 	t.Run("getOutAmountInBToken- Partial fill", func(t *testing.T) {
// 		fract := 0.501
// 		res, err := getOutAmountInBToken(ctx, newBids(), decimal.NewFromFloat(fract))
// 		assert.Equal(t, err, nil)
// 		assert.Equal(t, res.Size.String(), decimal.NewFromFloat(900*fract).String())
// 	})

// 	t.Run("getOutAmountInBToken- liquidity insufficient", func(t *testing.T) {
// 		_, err := getOutAmountInBToken(ctx, newBids(), decimal.NewFromFloat(1+2+3+1))
// 		assert.Equal(t, err, models.ErrInsufficientLiquity)
// 	})
// }
