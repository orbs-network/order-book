package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) GetQuote(ctx context.Context, symbol models.Symbol, side models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal) (models.QuoteRes, error) {

	// make sure inAmount is positivr
	if !inAmount.IsPositive() {
		return models.QuoteRes{}, models.ErrInAmount
	}
	var it models.OrderIter
	var res models.QuoteRes
	var err error
	if side == models.BUY {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		if it == nil {
			logctx.Error(ctx, "GetMinAsk failed")
			return models.QuoteRes{}, models.ErrIterFail
		}
		res, err = getOutAmountInAToken(ctx, it, inAmount)

	} else { // SELL
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		if it == nil {
			logctx.Error(ctx, "GetMaxBid failed")
			return models.QuoteRes{}, models.ErrIterFail
		}
		res, err = getOutAmountInBToken(ctx, it, inAmount)
	}
	if err != nil {
		logctx.Error(ctx, "getQuoteResIn failed", logger.Error(err))
		return models.QuoteRes{}, err
	}

	// apply min amount out threshold
	if minOutAmount != nil && (*minOutAmount).GreaterThanOrEqual(res.Size) {
		return models.QuoteRes{}, models.ErrMinOutAmount
	}

	return res, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getOutAmountInAToken(ctx context.Context, it models.OrderIter, inAmountB decimal.Decimal) (models.QuoteRes, error) {
	outAmountA := decimal.NewFromInt(0)
	var frags []models.OrderFrag
	var order *models.Order
	for it.HasNext() && inAmountB.IsPositive() {
		order = it.Next(ctx)
		// skip orders with locked funds
		if order.GetAvailableSize().IsPositive() {
			// max Spend in B token  for this order
			orderSizeB := order.Price.Mul(order.GetAvailableSize())
			// spend the min of orderSizeB/inAmountB
			spendB := decimal.Min(orderSizeB, inAmountB)

			// Gain
			gainA := spendB.Div(order.Price)

			// sub-add
			inAmountB = inAmountB.Sub(spendB)
			outAmountA = outAmountA.Add(gainA)

			// res
			logctx.Info(ctx, fmt.Sprintf("append OrderFrag gainA: %s", gainA.String()))
			logctx.Info(ctx, fmt.Sprintf("append OrderFrag spendB: %s", spendB.String()))
			frags = append(frags, models.OrderFrag{OrderId: order.Id, Size: gainA})
		}
	}
	// not all is Spent - error
	if inAmountB.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append OrderFrag outAmountA: %s", outAmountA.String()))
	return models.QuoteRes{Size: outAmountA, OrderFrags: frags}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getOutAmountInBToken(ctx context.Context, it models.OrderIter, inAmountA decimal.Decimal) (models.QuoteRes, error) {
	outAmountB := decimal.NewFromInt(0)
	var order *models.Order
	var frags []models.OrderFrag
	for it.HasNext() && inAmountA.IsPositive() {
		order = it.Next(ctx)

		// Spend
		spendA := decimal.Min(order.GetAvailableSize(), inAmountA)
		fmt.Println("sizeA ", spendA.String())

		// Gain
		gainB := order.Price.Mul(spendA)
		fmt.Println("gainB ", gainB.String())

		// sub-add
		inAmountA = inAmountA.Sub(spendA)
		outAmountB = outAmountB.Add(gainB)

		// res
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag spendA: %s", spendA.String()))
		logctx.Info(ctx, fmt.Sprintf("append OrderFrag gainB: %s", gainB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, Size: spendA})
	}
	if inAmountA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Info(ctx, fmt.Sprintf("append OrderFrag outAmountB: %s", outAmountB.String()))
	return models.QuoteRes{Size: outAmountB, OrderFrags: frags}, nil
}
