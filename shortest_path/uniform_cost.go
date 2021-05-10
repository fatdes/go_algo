package shortest_path

import (
	"container/heap"
)

type Node struct {
	Cost int
	Path []Vertex

	Vertex Vertex
}

type Vertex interface {
	Edges() []Edge
}

type Edge interface {
	Cost() int
	From() Vertex
	To() Vertex
}

func Find(from Vertex, to Vertex) (found bool, result *Node) {
	if from == nil || to == nil {
		return false, nil
	}

	if from == to {
		return true, &Node{
			Vertex: from,
			Cost:   0,
			Path: []Vertex{
				from,
			},
		}
	}

	pq := make(PriorityQueue, 1)
	pq[0] = NewInitialItem(
		&Node{
			Vertex: from,
			Cost:   0,
			Path: []Vertex{
				from,
			},
		},
		0,
	)
	heap.Init(&pq)

	explored := map[Vertex]bool{}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		node := item.value

		if node.Vertex == to {
			return true, node
		}

		explored[node.Vertex] = true

		for _, edge := range node.Vertex.Edges() {
			to := edge.To()
			if _, found := explored[to]; !found {
				path := make([]Vertex, len(node.Path)+1)
				copy(path, node.Path)
				path[len(path)-1] = to
				heap.Push(&pq, NewItem(&Node{
					Vertex: to,
					Cost:   node.Cost + edge.Cost(),
					Path:   path,
				}))
			}
		}
	}

	return false, nil
}
