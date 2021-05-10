// Modified from https://golang.org/pkg/container/heap/

package shortest_path

type Item struct {
	value *Node
	index int
}

func (item *Item) Value() *Node {
	return item.value
}

type PriorityQueue []*Item

func NewItem(node *Node) *Item {
	return &Item{
		value: node,
	}
}

func NewInitialItem(node *Node, index int) *Item {
	return &Item{
		value: node,
		index: index,
	}
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Pop lowest value first
	return pq[i].value.Cost < pq[j].value.Cost
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
