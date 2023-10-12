package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/orbs-network/order-book/models"
	"github.com/orbs-network/order-book/utils/logger"
	"github.com/orbs-network/order-book/utils/logger/logctx"
	"github.com/shopspring/decimal"
)

// type PairOrderBook struct {
// 	repo memoryRepository.InMemoryRepository
// }

// func NewPairOrderBook(repo data.Repository) *PairOrderBook {
// 	return &PairOrderBook{repo: repo}
// }

func (s *Service) AddOrder(ctx context.Context, price decimal.Decimal, symbol models.Symbol, size decimal.Decimal) (models.Order, error) {
	id := uuid.New()
	order := models.Order{
		Id:        id,
		Price:     price,
		Symbol:    symbol,
		Size:      size,
		Signature: nil,
		Pending:   false,
	}

	if err := s.store.StoreOrder(order); err != nil {
		logctx.Error(ctx, "failed to add order", logger.Error(err))
		return models.Order{}, err
	}

	logctx.Info(ctx, "order added", logger.String("ID", order.Id.String()), logger.String("price", order.Price.String()), logger.String("size", order.Size.String()))
	return order, nil
}

func (s *Service) CancelOrder(orderId uuid.UUID) {
	// Additional business logic or validations can be placed here
	s.store.RemoveOrder(orderId)

}

func (s *Service) GetBestOffer() *models.Order {
	prices := s.store.GetAllPrices()
	if len(prices) == 0 {
		return nil // or handle this case as appropriate for your application
	}

	bestPrice := prices[0]
	orders := s.store.GetOrdersAtPrice(bestPrice)
	if len(orders) == 0 {
		return nil // or handle this as appropriate, though this case should theoretically not occur
	}

	return &orders[0]
}

func (s *Service) GetVolumeAtLimit(price decimal.Decimal) decimal.Decimal {
	// Additional business logic can be placed here
	orders := s.store.GetOrdersAtPrice(price)
	var totalVolume decimal.Decimal
	for _, order := range orders {
		totalVolume = totalVolume.Add(order.Size)
	}
	return totalVolume
}
