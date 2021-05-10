package lwwes_test

import (
	"fatdes/go_algo/lwwes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func depthFirstAction(current *lwwes.TreeNode, action func(current *lwwes.TreeNode, child *lwwes.TreeNode)) {
	for _, child := range current.Children {
		action(current, child)
		depthFirstAction(child, action)
	}
}

var simpleTree = func(graph *lwwes.Graph, now int64) {
	graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
	graph.AddNode(&lwwes.Node{ID: "node1_1"}, &lwwes.Node{ID: "node1"}, now)
	graph.AddNode(&lwwes.Node{ID: "node1_2"}, &lwwes.Node{ID: "node1"}, now)
	graph.AddNode(&lwwes.Node{ID: "node1_3"}, &lwwes.Node{ID: "node1"}, now)
	graph.AddNode(&lwwes.Node{ID: "node2"}, lwwes.RootNode, now)
	graph.AddNode(&lwwes.Node{ID: "node2_1"}, &lwwes.Node{ID: "node2"}, now)
	graph.AddNode(&lwwes.Node{ID: "node2_2"}, &lwwes.Node{ID: "node2"}, now)
	graph.AddNode(&lwwes.Node{ID: "node3"}, lwwes.RootNode, now)
	graph.AddNode(&lwwes.Node{ID: "node3_1"}, &lwwes.Node{ID: "node3"}, now)
	graph.AddNode(&lwwes.Node{ID: "node4"}, lwwes.RootNode, now)
}

func Test_Graph_QueryNode(t *testing.T) {
	type tc struct {
		name   string
		setup  func(graph *lwwes.Graph, now int64)
		find   []string
		expect bool
	}

	tcs := []*tc{
		{
			name:   "find root",
			setup:  simpleTree,
			find:   []string{lwwes.RootNode.ID},
			expect: true,
		},
		{
			name:   "find level 1 node",
			setup:  simpleTree,
			find:   []string{"node1", "node2", "node3", "node4"},
			expect: true,
		},
		{
			name:   "find level 2 node",
			setup:  simpleTree,
			find:   []string{"node1_1", "node1_2", "node1_3", "node2_1", "node2_2", "node3_1"},
			expect: true,
		},
		{
			name: "removed node before add = exists",
			setup: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.RemoveNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now-100)
			},
			find:   []string{"node1", "node1_1", "node1_2", "node1_3"},
			expect: true,
		},
		{
			name: "should not find removed node and subtree",
			setup: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.RemoveNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now+100)
			},
			find:   []string{"node1", "node1_1", "node1_2", "node1_3"},
			expect: false,
		},
		{
			name: "cannot find orphan node",
			setup: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.AddNode(&lwwes.Node{ID: "node999"}, &lwwes.Node{ID: "node888"}, now)
			},
			find:   []string{"node999", "node888"},
			expect: false,
		},
		{
			name: "cannot node with multiple parent",
			setup: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.AddNode(&lwwes.Node{ID: "node4"}, &lwwes.Node{ID: "node1"}, now)
			},
			find:   []string{"node4"},
			expect: false,
		},
	}

	for _, ti := range tcs {
		graph := lwwes.NewGraph()
		now := time.Now().UnixNano()

		ti.setup(graph, now)

		for _, id := range ti.find {
			node := graph.QueryNode(id)
			assert.Equal(t, ti.expect, node != nil, ti.name)
		}
	}
}

type pair struct {
	from string
	to   string
}

var generateIncLevelPairs = func(count int) []*pair {
	pairs := make([]*pair, count)
	pairs[0] = &pair{from: lwwes.RootNode.ID, to: fmt.Sprintf("node%d", 1)}
	for i := 1; i < count; i++ {
		pairs[i] = &pair{from: fmt.Sprintf("node%d", i), to: fmt.Sprintf("node%d", i+1)}
	}
	return pairs
}

var generateSameLevelPairs = func(count int) []*pair {
	pairs := make([]*pair, count-1)
	pairs[0] = &pair{from: lwwes.RootNode.ID, to: fmt.Sprintf("node%d", 1)}
	for i := 2; i < count; i++ {
		pairs[i-1] = &pair{from: fmt.Sprintf("node%d", 1), to: fmt.Sprintf("node%d", i)}
	}
	return pairs
}

func Test_Graph_QueryAllNode(t *testing.T) {
	type tc struct {
		name   string
		setup  func(graph *lwwes.Graph, now int64)
		expect []*pair
	}

	tcs := []*tc{
		{
			name:   "root node only",
			setup:  func(_ *lwwes.Graph, _ int64) {},
			expect: []*pair{},
		},
		{
			name: "one level nodes",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node2"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node3"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node4"}, lwwes.RootNode, now)
			},
			expect: []*pair{
				{from: lwwes.RootNode.ID, to: "node1"},
				{from: lwwes.RootNode.ID, to: "node2"},
				{from: lwwes.RootNode.ID, to: "node3"},
				{from: lwwes.RootNode.ID, to: "node4"},
			},
		},
		{
			name: "two level nodes",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node1_1"}, &lwwes.Node{ID: "node1"}, now)
				graph.AddNode(&lwwes.Node{ID: "node1_2"}, &lwwes.Node{ID: "node1"}, now)
				graph.AddNode(&lwwes.Node{ID: "node1_3"}, &lwwes.Node{ID: "node1"}, now)
				graph.AddNode(&lwwes.Node{ID: "node2"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node2_1"}, &lwwes.Node{ID: "node2"}, now)
				graph.AddNode(&lwwes.Node{ID: "node2_2"}, &lwwes.Node{ID: "node2"}, now)
				graph.AddNode(&lwwes.Node{ID: "node3"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node3_1"}, &lwwes.Node{ID: "node3"}, now)
				graph.AddNode(&lwwes.Node{ID: "node4"}, lwwes.RootNode, now)
			},
			expect: []*pair{
				{from: lwwes.RootNode.ID, to: "node1"},
				{from: "node1", to: "node1_1"},
				{from: "node1", to: "node1_2"},
				{from: "node1", to: "node1_3"},
				{from: lwwes.RootNode.ID, to: "node2"},
				{from: "node2", to: "node2_1"},
				{from: "node2", to: "node2_2"},
				{from: lwwes.RootNode.ID, to: "node3"},
				{from: "node3", to: "node3_1"},
				{from: lwwes.RootNode.ID, to: "node4"},
			},
		},
		{
			name: "n level nodes",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				for i := 1; i < 99; i++ {
					graph.AddNode(&lwwes.Node{ID: fmt.Sprintf("node%d", i+1)}, &lwwes.Node{ID: fmt.Sprintf("node%d", i)}, now)
				}
			},
			expect: generateIncLevelPairs(99),
		},
		{
			name: "two parents - removed",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node2"}, lwwes.RootNode, now)
				// two parents
				graph.AddNode(&lwwes.Node{ID: "node3"}, &lwwes.Node{ID: "node1"}, now)
				graph.AddNode(&lwwes.Node{ID: "node3"}, &lwwes.Node{ID: "node2"}, now)
			},
			expect: []*pair{
				{from: lwwes.RootNode.ID, to: "node1"},
				{from: lwwes.RootNode.ID, to: "node2"},
			},
		},
		{
			name: "removed node/edge do NOT count as parent",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "node2"}, lwwes.RootNode, now)
				// node2 -> node3
				graph.AddNode(&lwwes.Node{ID: "node3"}, &lwwes.Node{ID: "node2"}, now)
				// node3 is removed
				graph.RemoveNode(&lwwes.Node{ID: "node3"}, &lwwes.Node{ID: "node2"}, now+100)
				// node1 -> node3
				graph.AddNode(&lwwes.Node{ID: "node3"}, &lwwes.Node{ID: "node1"}, now+200)
			},
			expect: []*pair{
				{from: lwwes.RootNode.ID, to: "node1"},
				{from: "node1", to: "node3"},
				{from: lwwes.RootNode.ID, to: "node2"},
			},
		},
		{
			name: "n parents - removed",
			setup: func(graph *lwwes.Graph, now int64) {
				graph.AddNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now)
				for i := 2; i < 99; i++ {
					graph.AddNode(&lwwes.Node{ID: fmt.Sprintf("node%d", i)}, &lwwes.Node{ID: fmt.Sprintf("node%d", 1)}, now)
				}
			},
			expect: generateSameLevelPairs(99),
		},
		{
			name: "1 orphan node - removed",
			setup: func(graph *lwwes.Graph, now int64) {
				// orphan1 is orphan, so as orphan2
				graph.AddNode(&lwwes.Node{ID: "orphan2"}, &lwwes.Node{ID: "orphan1"}, now)
			},
			expect: []*pair{},
		},
		{
			name: "chained orphan node - both removed",
			setup: func(graph *lwwes.Graph, now int64) {
				// orphan1 is orphan, so as orphan2
				graph.AddNode(&lwwes.Node{ID: "orphan2"}, &lwwes.Node{ID: "orphan1"}, now)
				// orphan2 is orphan, so as orphan3
				graph.AddNode(&lwwes.Node{ID: "orphan3"}, &lwwes.Node{ID: "orphan2"}, now)
			},
			expect: []*pair{},
		},
		{
			name: "removed node becomes orphan node - removed",
			setup: func(graph *lwwes.Graph, now int64) {
				// orphan4 and orphan5 are orphans as remove is latest
				graph.AddNode(&lwwes.Node{ID: "orphan4"}, lwwes.RootNode, now)
				graph.AddNode(&lwwes.Node{ID: "orphan5"}, &lwwes.Node{ID: "orphan4"}, now)
				graph.RemoveNode(&lwwes.Node{ID: "orphan4"}, lwwes.RootNode, now+100)
			},
			expect: []*pair{},
		},
	}

	for _, ti := range tcs {
		graph := lwwes.NewGraph()
		now := time.Now().UnixNano()

		ti.setup(graph, now)

		root := graph.QueryAllNodes()
		pairs := []*pair{}
		depthFirstAction(root, func(current *lwwes.TreeNode, child *lwwes.TreeNode) {
			pairs = append(pairs, &pair{from: current.Node.ID, to: child.Node.ID})
		})

		assert.Equal(t, ti.expect, pairs, ti.name)
	}

}

func Test_Graph_Merge(t *testing.T) {
	type tc struct {
		name    string
		setup1  func(graph *lwwes.Graph, now int64)
		setup2  func(graph *lwwes.Graph, now int64)
		execute func(graph1, graph2 *lwwes.Graph)
		verify  func(ti *tc, graph1, graph2 *lwwes.Graph)
	}

	tcs := []*tc{
		{
			name:   "merge same is the same",
			setup1: simpleTree,
			setup2: simpleTree,
			execute: func(graph1, graph2 *lwwes.Graph) {
				graph1.Merge(graph2)
			},
			verify: func(ti *tc, graph1, _ *lwwes.Graph) {
				root := graph1.QueryAllNodes()
				pairs := []*pair{}
				depthFirstAction(root, func(current *lwwes.TreeNode, child *lwwes.TreeNode) {
					pairs = append(pairs, &pair{from: current.Node.ID, to: child.Node.ID})
				})

				assert.Equal(t, []*pair{
					{from: lwwes.RootNode.ID, to: "node1"},
					{from: "node1", to: "node1_1"},
					{from: "node1", to: "node1_2"},
					{from: "node1", to: "node1_3"},
					{from: lwwes.RootNode.ID, to: "node2"},
					{from: "node2", to: "node2_1"},
					{from: "node2", to: "node2_2"},
					{from: lwwes.RootNode.ID, to: "node3"},
					{from: "node3", to: "node3_1"},
					{from: lwwes.RootNode.ID, to: "node4"},
				}, pairs, ti.name)
			},
		},
		{
			name:   "merge removed node will remove it and subtree",
			setup1: simpleTree,
			setup2: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.RemoveNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now+100)
			},
			execute: func(graph1, graph2 *lwwes.Graph) {
				graph1.Merge(graph2)
			},
			verify: func(ti *tc, graph1, _ *lwwes.Graph) {
				root := graph1.QueryAllNodes()
				pairs := []*pair{}
				depthFirstAction(root, func(current *lwwes.TreeNode, child *lwwes.TreeNode) {
					pairs = append(pairs, &pair{from: current.Node.ID, to: child.Node.ID})
				})

				assert.Equal(t, []*pair{
					{from: lwwes.RootNode.ID, to: "node2"},
					{from: "node2", to: "node2_1"},
					{from: "node2", to: "node2_2"},
					{from: lwwes.RootNode.ID, to: "node3"},
					{from: "node3", to: "node3_1"},
					{from: lwwes.RootNode.ID, to: "node4"},
				}, pairs, ti.name)
			},
		},
		{
			name:   "merge removed node will remove it and subtree - commutative",
			setup1: simpleTree,
			setup2: func(graph *lwwes.Graph, now int64) {
				simpleTree(graph, now)
				graph.RemoveNode(&lwwes.Node{ID: "node1"}, lwwes.RootNode, now+100)
			},
			execute: func(graph1, graph2 *lwwes.Graph) {
				graph2.Merge(graph1)
			},
			verify: func(ti *tc, _, graph2 *lwwes.Graph) {
				root := graph2.QueryAllNodes()
				pairs := []*pair{}
				depthFirstAction(root, func(current *lwwes.TreeNode, child *lwwes.TreeNode) {
					pairs = append(pairs, &pair{from: current.Node.ID, to: child.Node.ID})
				})

				assert.Equal(t, []*pair{
					{from: lwwes.RootNode.ID, to: "node2"},
					{from: "node2", to: "node2_1"},
					{from: "node2", to: "node2_2"},
					{from: lwwes.RootNode.ID, to: "node3"},
					{from: "node3", to: "node3_1"},
					{from: lwwes.RootNode.ID, to: "node4"},
				}, pairs, ti.name)
			},
		},
	}

	for _, ti := range tcs {
		graph1 := lwwes.NewGraph()
		graph2 := lwwes.NewGraph()
		now := time.Now().UnixNano()

		ti.setup1(graph1, now)
		ti.setup2(graph2, now)

		ti.execute(graph1, graph2)

		ti.verify(ti, graph1, graph2)
	}
}
