package service

import (
	"context"

	"github.com/orbs-network/order-book/models"
)

func (s *Service) GetSymbols(ctx context.Context) ([]models.Symbol, error) {

	return models.GetAllSymbols(), nil
}
