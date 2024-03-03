package service

import (
	"context"
	"strconv"
	"time"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) startPeriodicChecks() {
	secSwapStarted := restutils.GetEnv("SEC_PERIODIC_INTERVAL", "10")
	sec, _ := strconv.Atoi(secSwapStarted)
	ctx := context.Background()
	logctx.Info(ctx, "startPeriodicChecks", logger.Int("sec_interval", sec))
	go func() {
		interval := time.Tick(time.Second * time.Duration(sec))
		for range interval {
			s.periodicCheck(ctx)
		}
	}()
}
func (s *Service) periodicCheck(ctx context.Context) {
	// cleanup dangeling swaps which did not start
	secSwapStarted := restutils.GetEnv("SEC_SWAP_STARTED", "60")
	sec, _ := strconv.Atoi(secSwapStarted)

	if sec > 0 { // USE ZERO as Turn Off feature flag
		logctx.Info(ctx, "SEC_SWAP_STARTED", logger.Int("sec_interval", sec))
		err := s.checkNonStartedSwaps(ctx, int64(sec))
		if err != nil {
			logctx.Error(ctx, "Error in peridic checks", logger.Error(err))
		}
	}
}

func secondsSinceTimestamp(t time.Time) (int64, error) {

	// Calculate the duration since the parsed time
	duration := time.Since(t)

	// Extract the number of seconds from the duration
	seconds := int64(duration.Seconds())

	return seconds, nil
}

func (s *Service) checkSwapStarted(ctx context.Context, swap models.Swap, secPeriod int64) {
	// already started
	if swap.IsStarted() {
		return
	}
	// check if not started during a period of x sec
	sec, err := secondsSinceTimestamp(swap.Created)
	if err != nil {
		logctx.Error(ctx, "secondsSinceTimestamp failed", logger.String("created", swap.Created.String()), logger.Error(err))
		return
	}
	if sec < secPeriod {
		// no need to abort
		return
	}
	logctx.Info(ctx, "swap was not started after allowed period", logger.String("swapId", swap.Id.String()))
	err = s.AbortSwap(ctx, swap.Id)
	if err != nil {
		logctx.Error(ctx, "failed to AutoabortSwap", logger.String("created", swap.Created.String()), logger.Error(err))
		return
	}
	// SUCCESS
	logctx.Info(ctx, "Auto abortSwap after interval", logger.String("swapId", swap.Id.String()), logger.Int("secPerios", int(secPeriod)))
}

func (s *Service) checkNonStartedSwaps(ctx context.Context, secPeriod int64) error {
	openSwaps, err := s.orderBookStore.GetOpenSwaps(ctx)
	if err != nil {
		return err
	}
	for _, swap := range openSwaps {
		s.checkSwapStarted(ctx, swap, secPeriod)
	}
	return nil
}
