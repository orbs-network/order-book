package rest

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/orbs-network/order-book/models"
	"github.com/shopspring/decimal"
)

// Service represents the methods available on the service to handle the actual request.
type Service interface {
	AddOrder(ctx context.Context, price decimal.Decimal, symbol models.Symbol, size decimal.Decimal) (models.Order, error)
}

type Handler struct {
	svc    Service
	router *mux.Router
}

func NewHandler(svc Service, r *mux.Router) (*Handler, error) {
	if svc == nil {
		return nil, fmt.Errorf("svc cannot be nil")
	}

	if r == nil {
		return nil, fmt.Errorf("router cannot be nil")
	}

	return &Handler{
		svc:    svc,
		router: r,
	}, nil
}
