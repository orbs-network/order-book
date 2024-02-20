package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func (s *Service) startPeriodicChecks() {
	go func() {
		ctx := context.Background()
		interval := time.Tick(10 * time.Second)
		for range interval {
			s.periodicCheck(ctx)
		}
	}()
}
func (s *Service) periodicCheck(ctx context.Context) {
	// cleanup dangeling swaps which did not start
	secSwapStarted := restutils.GetEnv("SEC_SWAP_STARTED", "60")
	sec, _ := strconv.Atoi(secSwapStarted)
	err := s.checkNonStartedSwaps(ctx, int64(sec))
	if err != nil {
		logctx.Error(ctx, "Error in peridic checks", logger.Error(err))
	}
}

func secondsSinceTimestamp(t time.Time) (int64, error) {

	// Calculate the duration since the parsed time
	duration := time.Since(t)

	// Extract the number of seconds from the duration
	seconds := int64(duration.Seconds())

	return seconds, nil
}

func (s *Service) checkNonStartedSwaps(ctx context.Context, secPeriod int64) error {
	swapKeys, err := s.orderBookStore.EnumSubKeysOf(ctx, "swapId")
	if err != nil {
		return err
	}
	fmt.Println(swapKeys)
	for _, swapKey := range swapKeys {
		splt := strings.Split(swapKey, ":")
		if len(splt) > 0 {
			swapId := splt[1]
			uid, err := uuid.Parse(swapId)
			if err == nil {
				swap, err := s.orderBookStore.GetSwap(ctx, uid)
				if err != nil {
					logctx.Error(ctx, "Error swap not found", logger.String("swapid", swapId), logger.Error(err))
				} else {

					// check if not started
					sec, err := secondsSinceTimestamp(swap.Created)
					if err == nil {
						if sec > secPeriod {
							fmt.Println(swap.Created)
							fmt.Println("swap was not started- removing")
							s.AbortSwap(ctx, uid)
						}
					} else {
						logctx.Error(ctx, "secondsSinceTimestamp failed", logger.String("created", swap.Created.String()), logger.Error(err))
					}
				}
			} else {
				logctx.Error(ctx, "failed to parse swapId", logger.String("swapid", swapId), logger.Error(err))
			}
		} else {
			logctx.Error(ctx, "swapKey couldnt split on colone", logger.String("swapKey", swapKey))
		}
	}
	return nil
}
