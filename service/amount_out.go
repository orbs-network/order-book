package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) GetQuote(ctx context.Context, symbol models.Symbol, makerSide models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal, inDec, outDec int) (models.QuoteRes, error) {

	logctx.Info(ctx, "GetQuote started", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("inAmount", inAmount.String()))
	if minOutAmount != nil {
		logctx.Info(ctx, "GetQuote minOutAmount requested", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("minOutAmount", minOutAmount.String()))
	}

	// make sure inAmount is positivr
	if !inAmount.IsPositive() {
		return models.QuoteRes{}, models.ErrInAmount
	}
	var it models.OrderIter
	var res models.QuoteRes
	var err error
	if makerSide == models.SELL {
		it = s.orderBookStore.GetMinAsk(ctx, symbol)
		if it == nil {
			logctx.Error(ctx, "GetMinAsk failed")
			return models.QuoteRes{}, models.ErrIterFail
		}
		if !it.HasNext() {
			logctx.Warn(ctx, "insufficient liquidity", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("inAmount", inAmount.String()))
			return models.QuoteRes{}, models.ErrInsufficientLiquity
		}

		res, err = getOutAmountInAToken(ctx, it, inAmount, inDec, outDec)

	} else { // BUY
		it = s.orderBookStore.GetMaxBid(ctx, symbol)
		if it == nil {
			logctx.Warn(ctx, "GetMaxBid failed no orders in iterator")
			return models.QuoteRes{}, models.ErrIterFail
		}
		if !it.HasNext() {
			logctx.Warn(ctx, "GetMaxBid failed no orders in iterator")
			return models.QuoteRes{}, models.ErrInsufficientLiquity
		}
		res, err = getOutAmountInBToken(ctx, it, inAmount, inDec, outDec)
	}
	if err != nil {
		logctx.Error(ctx, "getQuoteResIn failed", logger.Error(err))
		return models.QuoteRes{}, err
	}

	// apply min amount out threshold
	if minOutAmount != nil {
		logctx.Info(ctx, "minOutAmount check", logger.String("minOutAmount", minOutAmount.String()), logger.String("amountOut", res.Size.String()))
		if minOutAmount.GreaterThan(res.Size) {
			logctx.Info(ctx, "minOutAmount was applied")
			return models.QuoteRes{}, models.ErrMinOutAmount
		}
	}

	logctx.Info(ctx, "GetQuote Finished OK", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("inAmount", inAmount.String()))

	return res, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getOutAmountInAToken(ctx context.Context, it models.OrderIter, inAmountB decimal.Decimal, inDec, outDec int) (models.QuoteRes, error) {
	outAmountA := decimal.NewFromInt(0)
	var frags []models.OrderFrag
	var order *models.Order
	for it.HasNext() && inAmountB.IsPositive() {
		order = it.Next(ctx)
		if order == nil {
			logctx.Error(ctx, "order::it.Next() returned nil")
			return models.QuoteRes{}, models.ErrUnexpectedError
		}
		// Unexpected to get cancelled orders in price list
		if order.Cancelled {
			logctx.Error(ctx, "cancelled order exists in the price list (ignore and continue)", logger.String("orderId", order.Id.String()))
			return models.QuoteRes{}, models.ErrUnexpectedError
		}
		// skip orders with locked funds
		if order.GetAvailableSize().IsPositive() {
			// max Spend in B token for this order
			orderSizeB := order.Price.Mul(order.GetAvailableSize())
			// spend the min of orderSizeB/inAmountB
			spendB := decimal.Min(orderSizeB, inAmountB)

			//Gain
			gainA := spendB.Div(order.Price)
			println("gainA ", gainA.String())

			//sub - add
			inAmountB = inAmountB.Sub(spendB)
			outAmountA = outAmountA.Add(gainA)

			// res
			logctx.Debug(ctx, fmt.Sprintf("Price: %s", order.Price.String()))
			logctx.Debug(ctx, fmt.Sprintf("append OrderFrag gainA: %s", gainA.String()))
			logctx.Debug(ctx, fmt.Sprintf("append OrderFrag spendB: %s", spendB.String()))
			frags = append(frags, models.OrderFrag{OrderId: order.Id, OutSize: gainA, InSize: spendB})
		}
	}
	// not all is Spent - error
	if inAmountB.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Debug(ctx, fmt.Sprintf("append OrderFrag outAmountA: %s", outAmountA.String()))
	return models.QuoteRes{Size: outAmountA, OrderFrags: frags}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getOutAmountInBToken(ctx context.Context, it models.OrderIter, inAmountA decimal.Decimal, inDec, outDec int) (models.QuoteRes, error) {
	outAmountB := decimal.NewFromInt(0)
	var order *models.Order
	var frags []models.OrderFrag
	for it.HasNext() && inAmountA.IsPositive() {
		order = it.Next(ctx)
		if order == nil {
			logctx.Error(ctx, "order::it.Next() returned nil")
			return models.QuoteRes{}, models.ErrUnexpectedError
		}
		// Unexpected to get cancelled orders in price list
		if order.Cancelled {
			logctx.Error(ctx, "order::it.Next() returned a cencelled order", logger.String("orderId", order.Id.String()))
			return models.QuoteRes{}, models.ErrUnexpectedError
		}
		// skip orders with locked funds
		if order.GetAvailableSize().IsPositive() {
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
			logctx.Debug(ctx, fmt.Sprintf("Price: %s", order.Price.String()))
			//logctx.Debug(ctx, fmt.Sprintf("Onchain Price: %s", ocPrice.String()))
			logctx.Debug(ctx, fmt.Sprintf("append OrderFrag spendA: %s", spendA.String()))
			logctx.Debug(ctx, fmt.Sprintf("append OrderFrag gainB: %s", gainB.String()))
			frags = append(frags, models.OrderFrag{OrderId: order.Id, OutSize: gainB, InSize: spendA})
		}
	}
	if inAmountA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Debug(ctx, fmt.Sprintf("append OrderFrag outAmountB: %s", outAmountB.String()))
	return models.QuoteRes{Size: outAmountB, OrderFrags: frags}, nil
}
