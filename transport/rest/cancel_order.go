package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/service"
	"github.com/orbs-network/order-book/utils"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
)

type CancelOrderResponse struct {
	OrderId string `json:"orderId"`
}

func (h *Handler) CancelOrderByOrderId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	orderIdStr := vars["orderId"]

	orderId, err := uuid.Parse(orderIdStr)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	logctx.Info(ctx, "user trying to cancel order by orderID", logger.String("userId", user.Id.String()), logger.String("orderId", orderId.String()))

	h.handleCancelOrder(hInput{
		ctx:         ctx,
		id:          orderId,
		isClientOId: false,
		userId:      user.Id,
		w:           w,
	})
}

func (h *Handler) CancelOrderByClientOId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := utils.GetUserCtx(ctx)
	if user == nil {
		logctx.Error(ctx, "user should be in context")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	clientOIdStr := vars["clientOId"]

	clientOId, err := uuid.Parse(clientOIdStr)
	if err != nil {
		http.Error(w, "invalid clientOId", http.StatusBadRequest)
		return
	}

	logctx.Info(ctx, "user trying to cancel order by clientOId", logger.String("userId", user.Id.String()), logger.String("clientOId", clientOId.String()))

	h.handleCancelOrder(hInput{
		ctx:         ctx,
		id:          clientOId,
		isClientOId: true,
		userId:      user.Id,
		w:           w,
	})
}

type hInput struct {
	ctx         context.Context
	id          uuid.UUID
	isClientOId bool
	userId      uuid.UUID
	w           http.ResponseWriter
}

// handleCancelOrder calls the service to cancel an order and writes the response to the client
func (h *Handler) handleCancelOrder(input hInput) {

	cancelledOrderId, err := h.svc.CancelOrder(input.ctx, service.CancelOrderInput{
		Id:          input.id,
		IsClientOId: input.isClientOId,
		UserId:      input.userId,
	})

	if err == models.ErrNotFound {
		logctx.Warn(input.ctx, "order not found", logger.String("id", input.id.String()))
		http.Error(input.w, "Order not found", http.StatusNotFound)
		return
	}

	if err == models.ErrUnauthorized {
		logctx.Warn(input.ctx, "user not authorized to cancel order", logger.String("id", input.id.String()))
		http.Error(input.w, "Not authorized", http.StatusUnauthorized)
		return
	}

	if err == models.ErrOrderPending {
		logctx.Warn(input.ctx, "cancelling order not possible when order is pending", logger.String("id", input.id.String()))
		http.Error(input.w, "Cannot cancel order due to pending fill", http.StatusConflict)
		return
	}

	if err == models.ErrOrderFilled {
		logctx.Warn(input.ctx, "cancelling order not possible when order is filled", logger.String("id", input.id.String()))
		http.Error(input.w, "Cannot cancel filled order", http.StatusConflict)
		return
	}

	if err != nil {
		logctx.Error(input.ctx, "failed to cancel order", logger.Error(err))
		http.Error(input.w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	if cancelledOrderId == nil {
		logctx.Error(input.ctx, "cancelled order ID is nil")
		http.Error(input.w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	res := CancelOrderResponse{
		OrderId: cancelledOrderId.String(),
	}

	resp, err := json.Marshal(res)
	if err != nil {
		logctx.Error(input.ctx, "failed to marshal created order", logger.Error(err))
		http.Error(input.w, "Error cancelling order. Try again later", http.StatusInternalServerError)
		return
	}

	input.w.Header().Set("Content-Type", "application/json")
	input.w.WriteHeader(http.StatusOK)

	if _, err := input.w.Write(resp); err != nil {
		logctx.Error(input.ctx, "failed to write response", logger.Error(err), logger.String("orderId", cancelledOrderId.String()))

		http.Error(input.w, "Error cancelling order. Try again later", http.StatusInternalServerError)
	}
}
