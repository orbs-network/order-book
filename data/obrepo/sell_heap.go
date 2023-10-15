package obrepo

// Min-heap for sell orders.
type sellHeap []*ordersAtPrice

func (h sellHeap) Len() int { return len(h) }
func (h sellHeap) Less(i, j int) bool {
	return h[i].Price.LessThan(h[j].Price) // Ascending order for sell orders (min heap)
}
func (h sellHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *sellHeap) Push(x interface{}) {
	n := len(*h)
	order := x.(*ordersAtPrice)
	order.Index = n
	*h = append(*h, order)
}

func (h *sellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	x.Index = -1
	*h = old[0 : n-1]
	return x
}
