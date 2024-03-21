package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/transport/restutils"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

const NUM_OF_ORDERS_LIMIT = 10

type CreateOrdersRequest struct {
	Symbol string               `json:"symbol"`
	Orders []CreateOrderRequest `json:"orders"`
}

type CreateOrdersResponse struct {
	Symbol  string         `json:"symbol"`
	Created []models.Order `json:"created"`
	Status  int            `json:"status"`
	Msg     string         `json:"msg"`
}

// TODO: use goroutines to create orders in parallel

func (h *Handler) CreateOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		restutils.WriteJSONError(ctx, w, http.StatusUnauthorized, "User not found")
		return
	}

	var args CreateOrdersRequest
	err := json.NewDecoder(r.Body).Decode(&args)
	if err != nil {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if len(args.Orders) == 0 {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, "Orders list is empty. Ensure you include 'symbol' and 'orders'")
		return
	}

	if len(args.Orders) > NUM_OF_ORDERS_LIMIT {
		restutils.WriteJSONError(ctx, w, http.StatusBadRequest, fmt.Sprintf("Maximum %d orders allowed", NUM_OF_ORDERS_LIMIT))
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
			eip712Msg:     &order.Eip712Msg,
		}); err != nil {
			logctx.Warn(ctx, "failed to validate required fields", logger.Error(err), logger.String("userId", user.Id.String()))
			response.Status = http.StatusBadRequest
			response.Msg = err.Error()
			response.Created = createdOrders
			restutils.WriteJSONResponse(ctx, w, http.StatusBadRequest, response, logger.String("userId", user.Id.String()))
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
			response.Status = http.StatusBadRequest
			response.Msg = err.Error()
			response.Created = createdOrders
			break
		}

		abiFragment, err := restutils.ConvertToAbiFragment(order.Eip712Msg)
		if err != nil {
			logctx.Warn(ctx, "failed to convert eip712Msg to abi fragment", logger.Error(err))
			response.Status = http.StatusBadRequest
			response.Msg = fmt.Errorf("failed to parse eip712Msg: %w", err).Error()
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
			AbiFragment:   abiFragment,
		})

		if err == models.ErrClashingClientOrderId {
			logctx.Warn(ctx, "order with clientOrderId already exists", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = http.StatusConflict
			response.Msg = fmt.Sprintf("Order with clientOrderId %q already exists. You must first cancel this order", order.ClientOId.String())
			break
		}

		if err == models.ErrClashingOrderId {
			logctx.Warn(ctx, "order with orderId already exists", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = http.StatusConflict
			response.Msg = "Clashing order details. Please retry"
			break
		}

		if err != nil {
			logctx.Error(ctx, "failed to create order", logger.Error(err), logger.String("clientOrderId", parsedFields.clientOrderId.String()), logger.String("userId", user.Id.String()))
			response.Status = http.StatusInternalServerError
			response.Msg = fmt.Sprintf("Error creating order with clientOrderId %q. Try again later", order.ClientOId.String())
			break
		}

		logctx.Debug(ctx, "user created order", logger.String("userId", user.Id.String()), logger.String("price", parsedFields.roundedDecPrice.String()), logger.String("size", parsedFields.decSize.String()), logger.String("clientOrderId", parsedFields.clientOrderId.String()))
		createdOrders = append(createdOrders, order)
	}

	response.Symbol = args.Symbol
	response.Created = createdOrders

	if len(createdOrders) != len(args.Orders) {
		logctx.Warn(ctx, "not all orders were created", logger.String("userId", user.Id.String()), logger.Int("numOfOrders", len(createdOrders)), logger.Int("numOfOrdersRequested", len(args.Orders)))
		restutils.WriteJSONResponse(ctx, w, http.StatusBadRequest, response, logger.String("userId", user.Id.String()))
		return
	}

	logctx.Debug(ctx, "user created orders", logger.String("userId", user.Id.String()), logger.Int("numOfOrders", len(createdOrders)))
	response.Status = http.StatusCreated
	restutils.WriteJSONResponse(ctx, w, http.StatusCreated, response, logger.String("userId", user.Id.String()))

}
