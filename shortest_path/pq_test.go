package shortest_path_test

import (
	"container/heap"
	"fatdes/go_algo/shortest_path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createNode(cost int) *shortest_path.Node {
	return &shortest_path.Node{
		Cost: cost,
	}
}

func Test_PQ(t *testing.T) {
	type tc struct {
		name   string
		init   []*shortest_path.Node
		push   []*shortest_path.Node
		expect []*shortest_path.Node
	}

	tcs := []*tc{
		{
			name:   "no item",
			expect: []*shortest_path.Node{},
		},
		{
			name: "init 1 item",
			init: []*shortest_path.Node{
				createNode(1),
			},
			expect: []*shortest_path.Node{
				createNode(1),
			},
		},
		{
			name: "init 2 item in desc order",
			init: []*shortest_path.Node{
				createNode(999),
				createNode(1),
			},
			expect: []*shortest_path.Node{
				createNode(1),
				createNode(999),
			},
		},
		{
			name: "init 2 item in desc order, push item in between",
			init: []*shortest_path.Node{
				createNode(999),
				createNode(1),
			},
			push: []*shortest_path.Node{
				createNode(555),
			},
			expect: []*shortest_path.Node{
				createNode(1),
				createNode(555),
				createNode(999),
			},
		},
	}

	for _, tt := range tcs {
		pq := make(shortest_path.PriorityQueue, 0)
		for i, node := range tt.init {
			pq = append(pq, shortest_path.NewInitialItem(node, i))
		}
		heap.Init(&pq)

		for _, node := range tt.push {
			heap.Push(&pq, shortest_path.NewItem(node))
		}

		actual := make([]*shortest_path.Node, 0)
		for pq.Len() > 0 {
			item := heap.Pop(&pq).(*shortest_path.Item)
			actual = append(actual, item.Value())
		}

		assert.Equal(t, tt.expect, actual)
	}
}
