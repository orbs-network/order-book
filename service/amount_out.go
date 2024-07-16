package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (s *Service) GetQuote(ctx context.Context, symbol models.Symbol, makerSide models.Side, inAmount decimal.Decimal, minOutAmount *decimal.Decimal, makerInToken string) (models.QuoteRes, error) {
	logctx.Info(ctx, "GetQuote started", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("inAmount", inAmount.String()))
	if minOutAmount != nil {
		logctx.Info(ctx, "GetQuote minOutAmount requested", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("minOutAmount", minOutAmount.String()))
	}

	// make sure inAmount is positivr
	if !inAmount.IsPositive() {
		return models.QuoteRes{}, models.ErrInAmount
	}

	// to verify onchain balance
	walletVerifier := NewWalletVerifier(makerInToken)

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
		res, err = getOutAmountInAToken(ctx, it, inAmount, walletVerifier)

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
		res, err = getOutAmountInBToken(ctx, it, inAmount, walletVerifier)
	}
	if err != nil {
		logctx.Warn(ctx, "getQuoteResIn failed", logger.Error(err))
		return models.QuoteRes{}, err
	}

	// apply min amount out threshold
	if minOutAmount != nil {
		logctx.Info(ctx, "minOutAmount check", logger.String("symbol", symbol.String()), logger.String("minOutAmount", minOutAmount.String()), logger.String("amountOut", res.Size.String()))
		if minOutAmount.GreaterThan(res.Size) {
			logctx.Info(ctx, "minOutAmount was applied", logger.String("symbol", symbol.String()))
			return models.QuoteRes{}, models.ErrMinOutAmount
		}
	}

	if !walletVerifier.CheckAll(ctx, s.orderBookStore) {
		logctx.Error(ctx, "walletVerifier CheckAll return false", logger.String("makerInToken", makerInToken), logger.String("makerInAmount", res.Size.String()))
		return models.QuoteRes{}, models.ErrInsufficientBalance
	}

	// apply on-chain balance verification on maker's InToken (which is out amount)
	logctx.Info(ctx, "GetQuote Finished OK", logger.String("symbol", symbol.String()), logger.String("makerSide", makerSide.String()), logger.String("inAmount", inAmount.String()))

	return res, nil
}

func validateOrder(ctx context.Context, order *models.Order) bool {
	if order == nil {
		logctx.Error(ctx, "iter_Next returned nil")
		return false
	}
	// Unexpected to get cancelled orders in price list
	if order.Cancelled {
		logctx.Error(ctx, "cancelled order exists in the price list", logger.String("orderId", order.Id.String()))
		return false
	}
	// skip orders with locked funds
	return order.GetAvailableSize().IsPositive()
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in B token (USD)
// amount out A token (ETH)
func getOutAmountInAToken(ctx context.Context, it models.OrderIter, inAmountB decimal.Decimal, verifier *WalletVerifier) (models.QuoteRes, error) {
	outAmountA := decimal.NewFromInt(0)
	var frags []models.OrderFrag
	var order *models.Order

	for it.HasNext() && inAmountB.IsPositive() {
		order = it.Next(ctx)
		if validateOrder(ctx, order) {
			// max Spend in B token for this order
			orderSizeB := order.Price.Mul(order.GetAvailableSize())
			// user spends B the min of orderSizeB / inAmountB
			takerSpendB := decimal.Min(orderSizeB, inAmountB)

			//Gain
			takerGainA := takerSpendB.Div(order.Price)

			// to verify onChain the maker can spend A token the taker gains
			verifier.Add(order.Signature.AbiFragment.Info.Swapper.String(), takerGainA)

			//sub - add
			inAmountB = inAmountB.Sub(takerSpendB)
			outAmountA = outAmountA.Add(takerGainA)

			// append
			frags = append(frags, models.OrderFrag{OrderId: order.Id, OutSize: takerGainA, InSize: takerSpendB})
			logctx.Debug(ctx, "getOutAmountInAToken - append order frag", logger.String("takerGainA", takerGainA.String()), logger.String("takerSpendB", takerSpendB.String()))
		}
	}
	// not all is Spent - error
	if inAmountB.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Debug(ctx, "getOutAmountInAToken total", logger.String("inAmountB", inAmountB.String()), logger.String("outAmountA", outAmountA.String()))
	return models.QuoteRes{Size: outAmountA, OrderFrags: frags}, nil
}

// PAIR/SYMBOL A-B (ETH-USDC)
// amount in A token (ETH)
// amount out B token (USD)
func getOutAmountInBToken(ctx context.Context, it models.OrderIter, inAmountA decimal.Decimal, verifier *WalletVerifier) (models.QuoteRes, error) {
	outAmountB := decimal.NewFromInt(0)
	var order *models.Order
	var frags []models.OrderFrag
	for it.HasNext() && inAmountA.IsPositive() {
		order = it.Next(ctx)
		if validateOrder(ctx, order) {
			// user Spends A
			takerSpendA := decimal.Min(order.GetAvailableSize(), inAmountA)

			// user Gains B
			takerGainB := order.Price.Mul(takerSpendA)

			// to verify onChain maker has the B funds for the user to gain
			verifier.Add(order.Signature.AbiFragment.Info.Swapper.String(), takerGainB)

			// sub-add
			inAmountA = inAmountA.Sub(takerSpendA)
			outAmountB = outAmountB.Add(takerGainB)

			// res
			frags = append(frags, models.OrderFrag{OrderId: order.Id, OutSize: takerGainB, InSize: takerSpendA})
			logctx.Debug(ctx, "getOutAmountInBToken append order frag", logger.String("takerGainB", takerGainB.String()), logger.String("takerSpendA", takerSpendA.String()))
		}
	}
	if inAmountA.IsPositive() {
		logctx.Warn(ctx, models.ErrInsufficientLiquity.Error())
		return models.QuoteRes{}, models.ErrInsufficientLiquity
	}
	logctx.Debug(ctx, "getOutAmountInBToken total", logger.String("inAmountA", inAmountA.String()), logger.String("outAmountB", outAmountB.String()))
	return models.QuoteRes{Size: outAmountB, OrderFrags: frags}, nil
}
