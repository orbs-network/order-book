package evmrepo

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (e *evmRepository) callContract(ctx context.Context, msg ethereum.CallMsg) (string, error) {
	res, err := e.client.CallContract(ctx, msg, nil)
	if err != nil {
		logctx.Error(ctx, "Error calling contract: %v", logger.Error(err))
		return "", err
	}

	ret := hex.EncodeToString(res)
	logctx.Debug(ctx, "CallContract ", logger.String("ret", ret))
	return ret, nil
}

func (e *evmRepository) BalanceOf(ctx context.Context, token, adrs string) (*big.Int, error) {
	tokAdrs := common.HexToAddress(token)
	mkrAdrs := common.HexToAddress(adrs)

	packed, err := e.tokenABI.Pack("balanceOf", mkrAdrs)
	if err != nil {
		return nil, err
	}

	// get balance
	callMsg := ethereum.CallMsg{
		To:   &tokAdrs,
		Data: packed,
	}

	hex, err := e.callContract(ctx, callMsg)
	if err != nil {
		logctx.Error(ctx, "callContract failed on balanceOf", logger.String("tokAdrs", token))
		return nil, err
	}

	// convert from hex
	balance := new(big.Int)
	_, success := balance.SetString(hex, 16)

	if !success {
		logctx.Error(ctx, "SrtString failed on balanceOf", logger.String("hex", hex))
		return nil, nil
	}

	return balance, nil
}

func (e *evmRepository) TokenDecimals(ctx context.Context, token, adrs string) (int64, error) {
	tokAdrs := common.HexToAddress(token)

	packed, err := e.tokenABI.Pack("decimals")
	if err != nil {
		return -1, err
	}

	// get balance
	callMsg := ethereum.CallMsg{
		To:   &tokAdrs,
		Data: packed,
	}

	hex, err := e.callContract(ctx, callMsg)
	if err != nil {
		logctx.Error(ctx, "callContract failed on TokenDecimals", logger.String("tokAdrs", token))
		return -1, err
	}

	dcmls, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		logctx.Error(ctx, "strconv.ParseInt failed", logger.String("tokAdrs", token))
		return -1, err
	}

	return dcmls, nil
}
