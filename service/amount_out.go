package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// orderID->amount bought or sold in A token always

func (s *Service) GetAmountOut(ctx context.Context, auctionId uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {

	var it models.OrderIter
	var res models.AmountOut
	var err error
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		res, err = getAmountOutInAToken(ctx, it, amountIn)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		res, err = getAmountOutInBToken(ctx, it, amountIn)
	}
	if err != nil {
		logctx.Error(ctx, "getAmountOutIn failed", logger.Error(err))
		return models.AmountOut{}, err
	}
	err = s.orderBookStore.StoreAuction(ctx, auctionId, res.OrderFrags)
	if err != nil {
		logctx.Error(ctx, "StoreAuction failed", logger.Error(err))
		return models.AmountOut{}, err
	}
	return res, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getAmountOutInAToken(ctx context.Context, it models.OrderIter, amountInB decimal.Decimal) (models.AmountOut, error) {
	amountOutA := decimal.NewFromInt(0)
	var frags []models.OrderFrag
	var order *models.Order
	for it.HasNext() && amountInB.IsPositive() {
		order = it.Next(ctx)
		// max Spend in B token  for this order
		orderSizeB := order.Price.Mul(order.Size)
		// spend the min of orderSizeB/amountInB
		spendB := decimal.Min(orderSizeB, amountInB)

		// Gain
		gainA := spendB.Div(order.Price)

		// sub-add
		amountInB = amountInB.Sub(spendB)
		logctx.Info(ctx, "StoreAuction failed")
		amountOutA = amountOutA.Add(gainA)

		// res
		logctx.Info(ctx, fmt.Sprintf("append FilledOrder gainA: %s", gainA.String()))
		logctx.Info(ctx, fmt.Sprintf("append FilledOrder spendB: %s", spendB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Amount: gainA})
	}
	// not all is Spent - error
	if amountInB.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append FilledOrder amountOutA: %s", amountOutA.String()))
	return models.AmountOut{AmountOut: amountOutA, OrderFrags: frags}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getAmountOutInBToken(ctx context.Context, it models.OrderIter, amountInA decimal.Decimal) (models.AmountOut, error) {
	amountOutB := decimal.NewFromInt(0)
	var order *models.Order
	var frags []models.OrderFrag
	for it.HasNext() && amountInA.IsPositive() {
		order = it.Next(ctx)

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
		logctx.Info(ctx, fmt.Sprintf("append FilledOrder spendA: %s", spendA.String()))
		logctx.Info(ctx, fmt.Sprintf("append FilledOrder gainB: %s", gainB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Amount: spendA})
	}
	if amountInA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append FilledOrder amountOutB: %s", amountOutB.String()))
	return models.AmountOut{AmountOut: amountOutB, OrderFrags: frags}, nil
}
