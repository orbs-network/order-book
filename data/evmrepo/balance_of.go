package evmrepo

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

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
	logctx.Debug(ctx, "Owner: %s", logger.String("ret", ret))
	return ret, nil
}

func (e *evmRepository) BalanceOf(ctx context.Context, token, adrs string) (*big.Int, error) {
	tokAdrs := common.HexToAddress(token)
	mkrAdrs := common.HexToAddress(adrs)
	fmt.Printf(tokAdrs.String(), mkrAdrs.String())

	packed, err := e.tokenABI.Pack("balanceOf", mkrAdrs)
	if err != nil {
		return nil, err
	}

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
