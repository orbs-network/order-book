package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type Reporter struct {
	stop        chan struct{}
	ctx         context.Context
	secInterval uint64
	svc         *Service
	fields      []logger.Field
}

func NewReporter(svc *Service) *Reporter {
	strSec := restutils.GetEnv("REPORT_SEC_INTERVAL", "10")
	num, err := strconv.ParseUint(strSec, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		num = 10
	}
	ctx := context.Background()
	logctx.Info(ctx, "NewReporter()", logger.String("SecInterval", strSec))

	return &Reporter{svc: svc, stop: make(chan struct{}), ctx: ctx, secInterval: num}
}

func (r *Reporter) logRoutine() {

	for {
		select {
		case <-r.stop:
			logctx.Info(r.ctx, "Reporter routine stopped.")
			return
		case <-time.After(time.Duration(r.secInterval) * time.Second):
			r.tick()
		}
	}
}

func (r *Reporter) Start() {
	logctx.Info(r.ctx, "Reporter logging routine.")
	go r.logRoutine()
}

func (r *Reporter) Stop() {
	logctx.Info(r.ctx, "Reporter logging routine.")
	close(r.stop)
}

func (r *Reporter) sumOrderSide(isAsk bool, it models.OrderIter) error {
	if !it.HasNext() {
		logctx.Warn(r.ctx, "GetMinAsk failed")
	}
	topOrder := float64(0)
	sumSize := float64(0)
	sumPending := float64(0)
	sumFilled := float64(0)
	openOrders := uint(0)

	topOrderName := "maxBid"
	side := "bid"
	if isAsk {
		side = "ask"
		topOrderName = "minAsk"
	}
	var order *models.Order
	for it.HasNext() {
		order = it.Next(r.ctx)
		// unexpected
		if order == nil {
			logctx.Error(r.ctx, "order::it.Next() returned nil")
			return models.ErrUnexpectedError
		}
		// sum & count
		if topOrder == 0 {
			// maxBid/minAsk only in first element
			topOrder = order.Price.InexactFloat64()
		}
		// agg
		sumSize += order.Size.InexactFloat64()
		sumPending += order.SizePending.InexactFloat64()
		sumFilled += order.SizeFilled.InexactFloat64()
		// count
		openOrders += 1

	}
	r.fields = append(r.fields,
		logger.Float64(topOrderName, topOrder),
		logger.Float64(side+"TotlSize", sumSize),
		logger.Float64(side+"TotlPendingSize", sumPending),
		logger.Float64(side+"TotlFilledSize", sumFilled),
		logger.Uint(side+"OpenOrders", openOrders),
	)
	return nil
}

func (r *Reporter) tick() {
	for _, sym := range models.GetAllSymbols() {
		// reset fields	per symbol
		r.fields = []logger.Field{logger.String("symbol", sym.String())}
		// ask
		itAsk := r.svc.orderBookStore.GetMinAsk(r.ctx, sym)
		if itAsk == nil {
			logctx.Error(r.ctx, "GetMinAsk failed")
			return
		}
		err := r.sumOrderSide(true, itAsk)
		if err != nil {
			logctx.Error(r.ctx, "sumOrderSide failed", logger.Error(err))
		}
		// bid
		itBid := r.svc.orderBookStore.GetMaxBid(r.ctx, sym)
		if itBid == nil {
			logctx.Error(r.ctx, "GetMinAsk failed")
			return
		}
		err = r.sumOrderSide(false, itBid)
		if err != nil {
			logctx.Error(r.ctx, "sumOrderSide failed", logger.Error(err))
		}
		// report
		logctx.Info(r.ctx, "report", r.fields...)
	}
}
