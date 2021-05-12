// Modified from https://golang.org/pkg/container/heap/

package shortest_path

type Item struct {
	value    interface{}
	priority PriorityFunc
	index    int
}

func (item *Item) Value() interface{} {
	return item.value
}

type PriorityQueue []*Item
type PriorityFunc func() int

func NewItem(value interface{}, priority PriorityFunc) *Item {
	return &Item{
		value:    value,
		priority: priority,
	}
}

func NewInitialItem(value interface{}, priority PriorityFunc, index int) *Item {
	return &Item{
		value:    value,
		priority: priority,
		index:    index,
	}
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Pop lowest value first
	return pq[i].priority() < pq[j].priority()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
