package memoryrepo

import "github.com/shopspring/decimal"

func (r *inMemoryRepository) GetAllPrices() []decimal.Decimal {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prices := make([]decimal.Decimal, len(r.sellOrders))
	i := 0
	for price := range r.sellOrders {
		prices[i] = decimal.RequireFromString(price)
		i++
	}
	return prices
}
