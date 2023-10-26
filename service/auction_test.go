package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/mocks"
	"github.com/orbs-network/order-book/service"
	"github.com/stretchr/testify/assert"
)

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

//	func newOrder(price, size int64) models.Order {
//		return models.Order{
//			Price: decimal.NewFromInt(price),
//			Size:  decimal.NewFromInt(size),
//		}
//	}
//
//	func newAsks() models.OrderIter {
//		return &iter{
//			orders: []models.Order{
//				newOrder(1000, 1),
//				newOrder(1001, 2),
//				newOrder(1002, 3),
//			},
//			index: -1,
//		}
//	}
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

// ///////////////////////////////////////////////////////////////
func TestService_ConfirmAuction(t *testing.T) {
	ctx := context.Background()
	fmt.Print(ctx)
	svc, _ := service.New(&mocks.MockOrderBookStore{Error: assert.AnError})
	//svc, _ := &mocks.MockOrderBookService{Error: models.ErrNoUserInContext},

	t.Run("ConfirmAuction- happy path", func(t *testing.T) {

		uuid, _ := uuid.NewUUID()
		res, err := svc.ConfirmAuction(ctx, uuid)
		fmt.Sprintf("%v", res)
		assert.Equal(t, err, nil)
		// assert.Equal(t, res.AmountOut.String(), decimal.NewFromFloat(1+2+3).String())
	})

}
