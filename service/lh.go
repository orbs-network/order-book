package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// orderID->amount bought or sold in A token always

func (s *Service) GetAmountOut(ctx context.Context, auctionId string, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {

	var it models.OrderIter
	var res models.AmountOut
	var err error
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		res, err = getAmountOutInAToken(it, amountIn)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		res, err = getAmountOutInBToken(it, amountIn)
	}
	if err != nil {
		return models.AmountOut{}, err
	}
	err = s.orderBookStore.StoreAuction(ctx, auctionId, res.FillOrders)
	if err != nil {
		return models.AmountOut{}, err
	}
	return res, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getAmountOutInAToken(it models.OrderIter, amountInB decimal.Decimal) (models.AmountOut, error) {
	amountOutA := decimal.NewFromInt(0)
	var fillOrders []models.FilledOrder
	var order *models.Order
	for it.HasNext() && amountInB.IsPositive() {
		order = it.Next()
		// max Spend in B token  for this order
		orderSizeB := order.Price.Mul(order.Size)
		// spend the min of orderSizeB/amountInB
		spendB := decimal.Min(orderSizeB, amountInB)

		// Gain
		gainA := spendB.Div(order.Price)

		// sub-add
		amountInB = amountInB.Sub(spendB)
		amountOutA = amountOutA.Add(gainA)

		// res
		fillOrders = append(fillOrders, models.FilledOrder{OrderId: order.Id, Amount: gainA})
	}
	// not all is Spent - error
	if amountInB.IsPositive() {
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}

	return models.AmountOut{AmountOut: amountOutA, FillOrders: fillOrders}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getAmountOutInBToken(it models.OrderIter, amountInA decimal.Decimal) (models.AmountOut, error) {
	amountOutB := decimal.NewFromInt(0)
	var order *models.Order
	var fillOrders []models.FilledOrder
	for it.HasNext() && amountInA.IsPositive() {
		order = it.Next()

		// Spend
		spendA := decimal.Min(order.Size, amountInA)
		fmt.Println("sizeA ", spendA.String())

		// Gain
		gainB := order.Price.Mul(spendA)
		fmt.Println("gainB ", gainB.String())

		// sub-add
		amountInA = amountInA.Sub(spendA)
		amountOutB = amountOutB.Add(gainB)

		// res
		fillOrders = append(fillOrders, models.FilledOrder{OrderId: order.Id, Amount: spendA})
	}
	if amountInA.IsPositive() {
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}

	return models.AmountOut{AmountOut: amountOutB, FillOrders: fillOrders}, nil
}
