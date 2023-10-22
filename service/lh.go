package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

type Order interface {
	GetMinPriceOrder() *models.Order
	GetNextOrder() *models.Order
}

// type fillStatus struct {
// 	OrderId string `json:"OrderId"`
// 	Fill    string `json:"Fill"`
// 	Symbol  string `json:"symbol"`
// 	Side    string `json:"side"`
// }

// orderID->amount bought or sold in A token always
type AmountOutRes struct {
	AmountOut  string            `json:"AmountOut"`
	FillOrders map[string]string `json:"FillOrders"`
}

func (s *Service) GetAmountOut(ctx context.Context, auctionID string, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (*AmountOutRes, error) {

	var it OrderIter
	var res *AmountOutRes
	var err error
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		res, err = getAmountOutInAToken(it, amountIn)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		res, err = getAmountOutInBToken(it, amountIn)
	}
	if err != nil {
		return nil, err
	}
	s.orderBookStore.StoreAuction(ctx, auctionID, res.FillOrders)
	return res, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getAmountOutInAToken(it OrderIter, amountInB decimal.Decimal) (*AmountOutRes, error) {
	amountOutA := decimal.NewFromInt(0)
	fillOrders := make(map[string]string)
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
		fillOrders[order.Id.String()] = gainA.String()
	}
	// not all is Spent - error
	if amountInB.IsPositive() {
		return nil, errors.New("not enough liquidity in book to sutisfy amountIn")
	}

	return &AmountOutRes{AmountOut: amountOutA.String(), FillOrders: fillOrders}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getAmountOutInBToken(it OrderIter, amountInA decimal.Decimal) (*AmountOutRes, error) {
	amountOutB := decimal.NewFromInt(0)
	var order *models.Order
	fillOrders := make(map[string]string)
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
		fillOrders[order.Id.String()] = spendA.String()
	}
	if amountInA.IsPositive() {
		return nil, errors.New("not enough liquidity in book to sutisfy amountIn")
	}
	//return amountOutB, nil
	return &AmountOutRes{AmountOut: amountOutB.String(), FillOrders: fillOrders}, nil
}
