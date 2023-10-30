package service

func (s *Service) GetStore() OrderBookStore {
	return s.orderBookStore
}
