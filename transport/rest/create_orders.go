package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

const NUM_OF_ORDERS_LIMIT = 10

type status string

const (
	SUCCESS status = "success"
	FAIL    status = "fail"
)

type CreateOrdersRequest struct {
	Symbol string               `json:"symbol"`
	Orders []CreateOrderRequest `json:"orders"`
}

type CreateOrdersResponse struct {
	Symbol        string         `json:"symbol"`
	Created       []models.Order `json:"created"`
	Status        status         `json:"status"`
	FailureReason string         `json:"failureReason"`
}

// TODO: use goroutines to create orders in parallel

func (h *Handler) CreateOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var args CreateOrdersRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if len(args.Orders) == 0 {
		http.Error(w, "Orders list is empty. Ensure you include 'symbol' and 'orders'", http.StatusBadRequest)
		return
	}

	if len(args.Orders) > NUM_OF_ORDERS_LIMIT {
		http.Error(w, fmt.Sprintf("Maximum %d orders allowed", NUM_OF_ORDERS_LIMIT), http.StatusBadRequest)
		return
	}

	var response CreateOrdersResponse
	response.Symbol = args.Symbol

	createdOrders := []models.Order{}

	for _, order := range args.Orders {
		if err = handleValidateRequiredFields(hVRFArgs{
			price:         order.Price,
			size:          order.Size,
			symbol:        args.Symbol,
			side:          order.Side,
			clientOrderId: order.ClientOrderId,
			eip712Sig:     order.Eip712Sig,
			eip712MsgData: &order.Eip712MsgData,
		}); err != nil {
			logctx.Warn(ctx, "failed to validate required fields", logger.Error(err), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = err.Error()
			response.Created = createdOrders
			writeJSONResponse(ctx, w, http.StatusBadRequest, response, logger.String("userId", user.Id.String()))
			return
		}
	}

	for _, order := range args.Orders {
		parsedFields, err := parseFields(w, pFInput{
			price:         order.Price,
			size:          order.Size,
			symbol:        args.Symbol,
			side:          order.Side,
			clientOrderId: order.ClientOrderId,
		})
		if err != nil {
			logctx.Warn(ctx, "failed to parse fields", logger.Error(err), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = err.Error()
			response.Created = createdOrders
			break
		}

		order, err := h.svc.CreateOrder(ctx, service.CreateOrderInput{
			UserId:        user.Id,
			Price:         parsedFields.roundedDecPrice,
			Symbol:        parsedFields.symbol,
			Size:          parsedFields.decSize,
			Side:          parsedFields.side,
			ClientOrderID: parsedFields.clientOrderId,
			Eip712Sig:     order.Eip712Sig,
			Eip712MsgData: order.Eip712MsgData,
		})

		if err == models.ErrSignatureVerificationError {
			logctx.Warn(ctx, "signature verification error", logger.Error(err), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = "Signature verification error"
			break
		}

		if err == models.ErrSignatureVerificationFailed {
			logctx.Warn(ctx, "signature verification failed", logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = "Signature verification failed"
			return
		}

		if err == models.ErrClashingClientOrderId {
			logctx.Warn(ctx, "order with clientOrderId already exists", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = fmt.Sprintf("Order with clientOrderId %q already exists. You must first cancel this order", order.ClientOId.String())
			break
		}

		if err == models.ErrClashingOrderId {
			logctx.Warn(ctx, "order with orderId already exists", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = "Clashing order details. Please retry"
			break
		}

		if err != nil {
			logctx.Error(ctx, "failed to create order", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = FAIL
			response.FailureReason = fmt.Sprintf("Error creating order with clientOrderId %q. Try again later", order.ClientOId.String())
			break
		}

		logctx.Info(ctx, "user created order", logger.String("userId", user.Id.String()), logger.String("price", parsedFields.roundedDecPrice.String()), logger.String("size", parsedFields.decSize.String()), logger.String("clientOrderId", parsedFields.clientOrderId.String()))
		createdOrders = append(createdOrders, order)
	}

	response.Symbol = args.Symbol
	response.Created = createdOrders

	if len(createdOrders) != len(args.Orders) {
		logctx.Warn(ctx, "not all orders were created", logger.String("userId", user.Id.String()), logger.Int("numOfOrders", len(createdOrders)), logger.Int("numOfOrdersRequested", len(args.Orders)))
		writeJSONResponse(ctx, w, http.StatusBadRequest, response, logger.String("userId", user.Id.String()))
		return
	}

	logctx.Info(ctx, "user created orders", logger.String("userId", user.Id.String()), logger.Int("numOfOrders", len(createdOrders)))
	response.Status = SUCCESS
	writeJSONResponse(ctx, w, http.StatusOK, response, logger.String("userId", user.Id.String()))

}
