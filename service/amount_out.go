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

func (s *Service) GetQuote(ctx context.Context, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {

	// make sure amountIn is positivr
	if !amountIn.IsPositive() {
		return models.AmountOut{}, models.ErrInAmount
	}
	var it models.OrderIter
	var res models.AmountOut
	var err error
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		if it == nil {
			logctx.Error(ctx, "GetMinAsk failed", logger.Error(err))
			return models.AmountOut{}, models.ErrIterFail
		}
		res, err = getAmountOutInAToken(ctx, it, amountIn)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		if it == nil {
			logctx.Error(ctx, "GetMaxBid failed", logger.Error(err))
			return models.AmountOut{}, models.ErrIterFail
		}
		res, err = getAmountOutInBToken(ctx, it, amountIn)
	}
	if err != nil {
		logctx.Error(ctx, "getAmountOutIn failed", logger.Error(err))
		return models.AmountOut{}, err
	}

	return res, nil
}

// orderID->amount bought or sold in A token always
func (s *Service) GetAmountOut(ctx context.Context, swapId uuid.UUID, symbol models.Symbol, side models.Side, amountIn decimal.Decimal) (models.AmountOut, error) {

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
	err = s.orderBookStore.StoreSwap(ctx, swapId, res.OrderFrags)
	if err != nil {
		logctx.Error(ctx, "StoreSwap failed", logger.Error(err))
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
		orderSizeB := order.Price.Mul(order.GetAvailableSize())
		// spend the min of orderSizeB/amountInB
		spendB := decimal.Min(orderSizeB, amountInB)

		// Gain
		gainA := spendB.Div(order.Price)

		// sub-add
		amountInB = amountInB.Sub(spendB)
		amountOutA = amountOutA.Add(gainA)

		// res
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag gainA: %s", gainA.String()))
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag spendB: %s", spendB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Size: gainA})
	}
	// not all is Spent - error
	if amountInB.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append OrderFrag amountOutA: %s", amountOutA.String()))
	return models.AmountOut{Size: amountOutA, OrderFrags: frags}, nil
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
		spendA := decimal.Min(order.GetAvailableSize(), amountInA)
		fmt.Println("sizeA ", spendA.String())

		// Gain
		gainB := order.Price.Mul(spendA)
		fmt.Println("gainB ", gainB.String())

		// sub-add
		amountInA = amountInA.Sub(spendA)
		amountOutB = amountOutB.Add(gainB)

		// res
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag spendA: %s", spendA.String()))
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag gainB: %s", gainB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Size: spendA})
	}
	if amountInA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.AmountOut{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append OrderFrag amountOutB: %s", amountOutB.String()))
	return models.AmountOut{Size: amountOutB, OrderFrags: frags}, nil
}
