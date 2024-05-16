package redisrepo

import (
	"context"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

func (r *redisRepository) GetMakerTokenBalance(ctx context.Context, token, wallet string) (decimal.Decimal, error) {
	val, err := r.ReadStrKey(ctx, GetMakerTokenTrackKey(token, wallet))
	res := decimal.NewFromInt(-1)
	if err != nil {
		logctx.Error(ctx, "GetMakerTokenBalance Failed to read key", logger.Error(err))
		return res, err
	}
	res, err = decimal.NewFromString(val)
	if err != nil {
		logctx.Error(ctx, "decimal NewFromString", logger.Error(err))
		return res, err
	}
	return res, nil
}
