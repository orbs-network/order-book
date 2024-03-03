package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (s *Service) checkSwapStarted(ctx context.Context, swapKey string, secPeriod int64) {
	splt := strings.Split(swapKey, ":")
	// key invalid
	if len(splt) < 2 {
		logctx.Error(ctx, "swapKey couldnt split on colone", logger.String("swapKey", swapKey))
		return
	}
	// parse swapID
	swapId := splt[2]
	uid, err := uuid.Parse(swapId)
	if err != nil {
		logctx.Error(ctx, "uuid failed to parse swapId", logger.String("swapid", swapId), logger.Error(err))
		return
	}
	// not found
	swap, err := s.orderBookStore.GetSwap(ctx, uid)
	if err != nil {
		logctx.Error(ctx, "Error swap not found", logger.String("swapId", swapId), logger.Error(err))
		return
	}

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
	logctx.Info(ctx, "swap was not started after allowed period", logger.String("swapId", swapId))
	err = s.AbortSwap(ctx, uid)
	if err != nil {
		logctx.Error(ctx, "failed to AutoabortSwap", logger.String("created", swap.Created.String()), logger.Error(err))
		return
	}
	// SUCCESS
	logctx.Info(ctx, "Auto abortSwap after interval", logger.String("swapId", swapId), logger.Int("secPerios", int(secPeriod)))
}

func (s *Service) checkNonStartedSwaps(ctx context.Context, secPeriod int64) error {
	swapKeys, err := s.orderBookStore.EnumSubKeysOf(ctx, "swap:open")
	if err != nil {
		return err
	}
	for _, swapKey := range swapKeys {
		s.checkSwapStarted(ctx, swapKey, secPeriod)
	}
	return nil
}
