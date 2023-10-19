package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type Order interface {
	GetMinPriceOrder() *models.Order
	GetNextOrder() *models.Order
}

func (s *Service) GetAmountOut(ctx context.Context, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (decimal.Decimal, error) {

	var it OrderIter
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		return getAmountOutInAToken(it, amountIn)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		return getAmountOutInBToken(it, amountIn)
	}
}

// func getAmountOutInAToken(it OrderIter, amountIn decimal.Decimal) (decimal.Decimal, error) {
// 	// buy 2 eth for 2000 usd
// 	// amount In = 2000, price = 1000
// 	amountOut := decimal.NewFromInt(0)
// 	var order *models.Order
// 	for it.HasNext() && amountIn.IsPositive() {
// 		fmt.Printf("it %v", &it)
// 		order = it.Next()

// 		fmt.Printf("orderPrice:\t", order.Price.String())
// 		fmt.Printf("orderSize:\t", order.Size.String())
// 		fmt.Printf(order.Size.String())

// 		// max buy
// 		maxBuySize := amountIn.Div(order.Price)
// 		fmt.Printf("maxBuySize Out")
// 		fmt.Printf(maxBuySize.String())

// 		minBuySize := decimal.Min(maxBuySize, order.Size)
// 		fmt.Printf("minBuySize Out")
// 		fmt.Printf(minBuySize.String())

// 		spent := minBuySize.Mul(order.Price)
// 		fmt.Printf("spent amount In")
// 		fmt.Printf(spent.String())

// 		amountIn = amountIn.Sub(spent)
// 		amountOut = amountOut.Add(minBuySize)
// 	}
// 	return amountOut, nil
// }

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getAmountOutInAToken(it OrderIter, amountInB decimal.Decimal) (decimal.Decimal, error) {
	amountOutA := decimal.NewFromInt(0)
	var order *models.Order
	for it.HasNext() && amountInB.IsPositive() {
		order = it.Next()
		fmt.Println("amountInB:\t", amountInB.String())
		fmt.Println("orderPrice:\t", order.Price.String())
		fmt.Println("orderSize:\t", order.Size.String())

		// max Spend in B token  for this order
		orderSizeB := order.Price.Mul(order.Size)
		fmt.Println("orderSizeB ", orderSizeB.String())

		// spend the min of orderSizeB/amountInB
		spendB := decimal.Min(orderSizeB, amountInB)
		fmt.Println("spendB ", spendB.String())
		// Gain
		gainA := spendB.Div(order.Price)
		fmt.Println("gainA ", gainA.String())

		amountInB = amountInB.Sub(spendB)
		fmt.Println("amountInB ", amountInB.String())
		amountOutA = amountOutA.Add(gainA)
		fmt.Println("amountOutA ", gainA.String())
	}
	return amountOutA, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getAmountOutInBToken(it OrderIter, amountInA decimal.Decimal) (decimal.Decimal, error) {
	amountOutB := decimal.NewFromInt(0)
	var order *models.Order
	for it.HasNext() && amountInA.IsPositive() {
		order = it.Next()

		fmt.Println("orderPrice:\t", order.Price.String())
		fmt.Println("orderSize:\t", order.Size.String())

		// Spend
		spendA := decimal.Min(order.Size, amountInA)
		fmt.Println("sizeA ", spendA.String())

		// Gain
		gainB := order.Price.Mul(spendA)
		fmt.Println("gainB ", gainB.String())

		amountInA = amountInA.Sub(spendA)
		amountOutB = amountOutB.Add(gainB)
	}
	return amountOutB, nil
}
