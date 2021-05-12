package shortest_path_test

import (
	"container/heap"
	"fatdes/go_algo/shortest_path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestItem struct {
	priority int
}

func (t *TestItem) Priority() int {
	return t.priority
}

func createTestItem(priority int) *TestItem {
	return &TestItem{
		priority: priority,
	}
}

func Test_PQ(t *testing.T) {
	type tc struct {
		name   string
		init   []*TestItem
		push   []*TestItem
		expect []*TestItem
	}

	tcs := []*tc{
		{
			name:   "no item",
			expect: []*TestItem{},
		},
		{
			name: "init 1 item",
			init: []*TestItem{
				createTestItem(1),
			},
			expect: []*TestItem{
				createTestItem(1),
			},
		},
		{
			name: "init 2 item in desc order",
			init: []*TestItem{
				createTestItem(999),
				createTestItem(1),
			},
			expect: []*TestItem{
				createTestItem(1),
				createTestItem(999),
			},
		},
		{
			name: "init 2 item in desc order, push item in between",
			init: []*TestItem{
				createTestItem(999),
				createTestItem(1),
			},
			push: []*TestItem{
				createTestItem(555),
			},
			expect: []*TestItem{
				createTestItem(1),
				createTestItem(555),
				createTestItem(999),
			},
		},
	}

	for _, tt := range tcs {
		pq := make(shortest_path.PriorityQueue, 0)
		for i, ti := range tt.init {
			pq = append(pq, shortest_path.NewInitialItem(ti, ti.Priority, i))
		}
		heap.Init(&pq)

		for _, ti := range tt.push {
			heap.Push(&pq, shortest_path.NewItem(ti, ti.Priority))
		}

		actual := make([]*TestItem, 0)
		for pq.Len() > 0 {
			item := heap.Pop(&pq).(*shortest_path.Item)
			actual = append(actual, item.Value().(*TestItem))
		}

		assert.Equal(t, tt.expect, actual)
	}
}
