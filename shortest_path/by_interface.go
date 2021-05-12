package shortest_path

import (
	"container/heap"
)

type Vertex interface {
	Edges() []Edge
}

type Edge interface {
	Cost() int
	From() Vertex
	To() Vertex
}

type node struct {
	vertex    Vertex
	totalCost int
	path      []interface{}
}

func (n *node) cost() int {
	return n.totalCost
}

type byInterface struct {
}

func NewUniformCostByInterface() UniformCost {
	return &byInterface{}
}

func (b *byInterface) Find(from interface{}, to interface{}) *Result {
	if from == nil || to == nil {
		return &Result{Found: false}
	}

	if from == to {
		return &Result{
			Found: true,
			Cost:  0,
			Path: []interface{}{
				from,
			},
		}
	}

	pq := make(PriorityQueue, 1)
	initialNode := &node{
		vertex:    from.(Vertex),
		totalCost: 0,
		path: []interface{}{
			from,
		},
	}
	pq[0] = NewInitialItem(
		initialNode,
		initialNode.cost,
		0,
	)
	heap.Init(&pq)

	explored := map[Vertex]bool{}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		n := item.value.(*node)

		if n.vertex == to {
			return &Result{
				Found: true,

				Cost: n.totalCost,
				Path: n.path,
			}
		}

		explored[n.vertex] = true

		for _, edge := range n.vertex.Edges() {
			to := edge.To()
			if _, found := explored[to]; !found {
				path := make([]interface{}, len(n.path)+1)
				copy(path, n.path)
				path[len(path)-1] = to
				newNode := &node{
					vertex:    to,
					totalCost: n.totalCost + edge.Cost(),
					path:      path,
				}
				heap.Push(&pq, NewItem(
					newNode,
					newNode.cost,
				))
			}
		}
	}

	return &Result{Found: false}
}
