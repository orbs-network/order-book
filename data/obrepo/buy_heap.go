package obrepo

// Max-heap for buy orders.
type buyHeap []*ordersAtPrice

func (h buyHeap) Len() int { return len(h) }
func (h buyHeap) Less(i, j int) bool {
	return h[i].Price.GreaterThan(h[j].Price) // Descending order for buy orders (max heap)
}
func (h buyHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *buyHeap) Push(x interface{}) {
	*h = append(*h, x.(*ordersAtPrice))
}

func (h *buyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
