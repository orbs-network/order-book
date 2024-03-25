package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

func str2time(timestamp string) (time.Time, error) {
	milliseconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	seconds := milliseconds / 1000
	nanoseconds := (milliseconds % 1000) * 1000000
	return time.Unix(seconds, nanoseconds), nil
}

// e.g 1648230000000
func getStartEndTime(r *http.Request) (time.Time, time.Time) {
	var startAt time.Time
	var endAt time.Time

	// end default NOW()
	qEndAt := r.URL.Query().Get("endAt")
	if qEndAt == "" {
		endAt = time.Now()
	} else {
		res, err := str2time(qEndAt)
		if err == nil {
			endAt = res
		} else {
			logctx.Warn(r.Context(), "fail to parse endAt", logger.String("endAt", qEndAt))
		}
	}
	// start
	qStartAt := r.URL.Query().Get("startAt")
	if qStartAt != "" {
		res, err := str2time(qStartAt)
		if err == nil {
			startAt = res
		} else {
			logctx.Warn(r.Context(), "fail to parse startAt", logger.String("startAt", qStartAt))
		}
	}

	// default 24h before endTime
	if startAt.IsZero() || startAt.After(endAt) {
		startAt = endAt.Add(-24 * time.Hour)
	}

	return startAt, endAt
}
