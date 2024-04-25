package service

import (
	"context"
	"fmt"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) GetQuote(ctx context.Context, symbol models.Symbol, side models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal, inDec, outDec int) (models.QuoteRes, error) {

	logctx.Info(ctx, "GetQuote started", logger.String("symbol", symbol.String()), logger.String("side", side.String()), logger.String("inAmount", inAmount.String()))
	if minOutAmount != nil {
		logctx.Info(ctx, "GetQuote minOutAmount requested", logger.String("symbol", symbol.String()), logger.String("side", side.String()), logger.String("minOutAmount", minOutAmount.String()))
	}

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
		if !it.HasNext() {
			logctx.Warn(ctx, "insufficient liquidity", logger.String("symbol", symbol.String()), logger.String("side", side.String()), logger.String("inAmount", inAmount.String()))
			return models.QuoteRes{}, models.ErrInsufficientLiquity
		}

		res, err = getOutAmountInAToken(ctx, it, inAmount, inDec, outDec)

	} else { // SELL
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

	logctx.Info(ctx, "GetQuote Finished OK", logger.String("symbol", symbol.String()), logger.String("side", side.String()), logger.String("inAmount", inAmount.String()))

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

		// skip orders with locked funds
		if order.Cancelled {
			logctx.Error(ctx, "cancelled order exists in the price list (ignore and continue)", logger.String("orderId", order.Id.String()))
			return models.QuoteRes{}, models.ErrUnexpectedError
		}
		// skip orders with locked funds or cancelled
		if order.GetAvailableSize().IsPositive() && !order.Cancelled {
			// calc onchain price to match solidity percision
			ocPrice, err := order.OnchainPrice(inDec, outDec)
			if err != nil {
				logctx.Error(ctx, "Onchain price failed for order", logger.String("orderId", order.Id.String()), logger.Error(err))
				return models.QuoteRes{}, models.ErrUnexpectedError
			}
			// max Spend in B token  for this order
			orderSizeB := ocPrice.Mul(order.GetAvailableSize())
			// spend the min of orderSizeB/inAmountB
			spendB := decimal.Min(orderSizeB, inAmountB)

			// Gain
			oldGainA := spendB.Div(order.Price)
			println("oldGainA ", oldGainA.String())

			// gain A onchain calc
			//fill = takerIn * orderIn / orderOut
			fmt.Println("orderIn ", order.Signature.AbiFragment.Input.Amount.String())
			fmt.Println("orderOut ", order.Signature.AbiFragment.Outputs[0].Amount.String())
			orderIn := decimal.NewFromBigInt(order.Signature.AbiFragment.Input.Amount, 0).Div(decimal.NewFromFloat(1e18))
			mulIn := spendB.Mul(orderIn)
			orderOut := decimal.NewFromBigInt(order.Signature.AbiFragment.Outputs[0].Amount, 0).Div(decimal.NewFromFloat(1e6))
			// round up
			gainA := mulIn.Div(orderOut).RoundUp(18)
			fmt.Println("gainA ", gainA.String())

			// sub-add
			inAmountB = inAmountB.Sub(spendB)
			outAmountA = outAmountA.Add(gainA)

			// res
			logctx.Debug(ctx, fmt.Sprintf("Price: %s", order.Price.String()))
			logctx.Debug(ctx, fmt.Sprintf("Onchain Price: %s", ocPrice.String()))
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

		// Spend
		spendA := decimal.Min(order.GetAvailableSize(), inAmountA)
		fmt.Println("sizeA ", spendA.String())

		// calc onchain price to match solidity percision
		// replace in and out decimals to match the order's side
		ocPrice, err := order.OnchainPrice(inDec, outDec)
		if err != nil {
			logctx.Error(ctx, "Onchain price failed for order", logger.String("orderId", order.Id.String()), logger.Error(err))
			return models.QuoteRes{}, models.ErrUnexpectedError
		}

		// Gain
		oldGainB := order.Price.Mul(spendA)
		fmt.Println("oldGainB ", oldGainB.String())

		// gain B onchain calc
		//fill = takerIn * orderIn / orderOut
		fmt.Println("orderIn ", order.Signature.AbiFragment.Input.Amount.String())
		fmt.Println("orderOut ", order.Signature.AbiFragment.Outputs[0].Amount.String())
		orderIn := decimal.NewFromBigInt(order.Signature.AbiFragment.Input.Amount, 0).Div(decimal.NewFromFloat(1e6))
		mulIn := spendA.Mul(orderIn)
		orderOut := decimal.NewFromBigInt(order.Signature.AbiFragment.Outputs[0].Amount, 0).Div(decimal.NewFromFloat(1e18))
		// Round UP
		gainB := mulIn.Div(orderOut).RoundUp(6)
		fmt.Println("gainB ", gainB.String())

		// sub-add
		inAmountA = inAmountA.Sub(spendA)
		outAmountB = outAmountB.Add(gainB)

		// res
		logctx.Debug(ctx, fmt.Sprintf("Price: %s", order.Price.String()))
		logctx.Debug(ctx, fmt.Sprintf("Onchain Price: %s", ocPrice.String()))
		logctx.Debug(ctx, fmt.Sprintf("append OrderFrag spendA: %s", spendA.String()))
		logctx.Debug(ctx, fmt.Sprintf("append OrderFrag gainB: %s", gainB.String()))
		frags = append(frags, models.OrderFrag{OrderId: order.Id, OutSize: gainB, InSize: spendA})
	}
	if inAmountA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Debug(ctx, fmt.Sprintf("append OrderFrag outAmountB: %s", outAmountB.String()))
	return models.QuoteRes{Size: outAmountB, OrderFrags: frags}, nil
}
