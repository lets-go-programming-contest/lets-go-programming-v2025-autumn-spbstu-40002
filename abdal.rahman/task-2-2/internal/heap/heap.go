package heap

type MaxHeap []int

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(val int) {
*h = append(*h, val)
}

func (h *MaxHeap) Pop() int {
if len(*h) == 0 {
panic("heap is empty")
}
res := (*h)[0]
*h = (*h)[1:]
return res
}
