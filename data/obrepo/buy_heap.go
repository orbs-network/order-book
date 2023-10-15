package obrepo

// Max-heap for buy orders.
type buyHeap []*ordersAtPrice

func (h buyHeap) Len() int { return len(h) }
func (h buyHeap) Less(i, j int) bool {
	return h[i].Price.GreaterThan(h[j].Price) // Descending order for buy orders (max heap)
}
func (h buyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *buyHeap) Push(x interface{}) {
	n := len(*h)
	order := x.(*ordersAtPrice)
	order.Index = n
	*h = append(*h, x.(*ordersAtPrice))
}

func (h *buyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	x.Index = -1 // Reset the index when popping from the heap
	*h = old[0 : n-1]
	return x
}


