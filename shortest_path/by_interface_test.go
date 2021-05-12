package shortest_path_test

import (
	"fatdes/go_algo/shortest_path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ByInterfaceString(vs []interface{}) string {
	ids := make([]string, len(vs))
	for i, v := range vs {
		ids[i] = v.(*testByInterfaceVertex).String()
	}

	return strings.Join(ids, ",")
}

type testByInterfaceVertex struct {
	id    string
	edges []shortest_path.Edge
}

func (v *testByInterfaceVertex) Edges() []shortest_path.Edge {
	return v.edges
}

func (v *testByInterfaceVertex) addEdge(to *testByInterfaceVertex, cost int) *testByInterfaceVertex {
	v.edges = append(v.edges, &testByInterfaceEdge{
		cost: cost, from: v, to: to,
	})
	return v
}

func (v *testByInterfaceVertex) String() string {
	return v.id
}

type testByInterfaceEdge struct {
	cost int
	from *testByInterfaceVertex
	to   *testByInterfaceVertex
}

func (e *testByInterfaceEdge) Cost() int {
	return e.cost
}

func (e *testByInterfaceEdge) From() shortest_path.Vertex {
	return e.from
}

func (e *testByInterfaceEdge) To() shortest_path.Vertex {
	return e.to
}

type testByInterfaceGraph struct {
	vs map[string]*testByInterfaceVertex
}

func (graph *testByInterfaceGraph) buildTestByInterfaceGraph() {
	graph.vs = map[string]*testByInterfaceVertex{
		"a": {id: "a"},
		"b": {id: "b"},
		"c": {id: "c"},
		"d": {id: "d"},
		"e": {id: "e"},
		"f": {id: "f"},
		"g": {id: "g"},
	}

	graph.vs["a"].addEdge(graph.vs["d"], 3).addEdge(graph.vs["b"], 5)
	graph.vs["b"].addEdge(graph.vs["c"], 1)
	graph.vs["c"].addEdge(graph.vs["e"], 6).addEdge(graph.vs["g"], 8)
	graph.vs["d"].addEdge(graph.vs["e"], 2).addEdge(graph.vs["f"], 2)
	graph.vs["e"].addEdge(graph.vs["b"], 4)
	graph.vs["f"].addEdge(graph.vs["g"], 3)
	graph.vs["g"].addEdge(graph.vs["e"], 4)
}

func Test_UniformCostByInterface_TestEmptyFrom(t *testing.T) {
	graph := &testByInterfaceGraph{}
	graph.buildTestByInterfaceGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(nil, graph.vs["b"])
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByInterface_TestEmptyTo(t *testing.T) {
	graph := &testByInterfaceGraph{}
	graph.buildTestByInterfaceGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], nil)
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByInterface_TestNotFound(t *testing.T) {
	graph := &testByInterfaceGraph{}
	graph.buildTestByInterfaceGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], &testByInterfaceVertex{id: "h"})
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByInterface_TestSearchSameNode(t *testing.T) {
	graph := &testByInterfaceGraph{}
	graph.buildTestByInterfaceGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], graph.vs["a"])
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 0)
	assert.Equal(t, ByInterfaceString(actual.Path), "a")
}

func Test_UniformCostByInterface_TestShortestPathFound(t *testing.T) {
	graph := &testByInterfaceGraph{}
	graph.buildTestByInterfaceGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], graph.vs["g"])
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 8)
	assert.Equal(t, ByInterfaceString(actual.Path), "a,d,f,g")
}
