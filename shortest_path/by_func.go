package shortest_path

import (
	"container/heap"
)

type edges func(interface{}) []interface{}
type edgeEnd func(interface{}) interface{}
type edgeCost func(interface{}) int

type byFunc struct {
	edges    edges
	edgeEnd  edgeEnd
	edgeCost edgeCost
}

func NewUniformCostByFunc(edges edges, edgeEnd edgeEnd, edgeCost edgeCost) *byFunc {
	return &byFunc{
		edges:    edges,
		edgeEnd:  edgeEnd,
		edgeCost: edgeCost,
	}
}

func (b *byFunc) Find(from interface{}, to interface{}) *Result {
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
		vertex:    from,
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

	explored := map[interface{}]bool{}

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

		for _, edge := range b.edges(n.vertex) {
			to := b.edgeEnd(edge)
			if _, found := explored[to]; !found {
				path := make([]interface{}, len(n.path)+1)
				copy(path, n.path)
				path[len(path)-1] = to
				newNode := &node{
					vertex:    to,
					totalCost: n.totalCost + b.edgeCost(edge),
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
