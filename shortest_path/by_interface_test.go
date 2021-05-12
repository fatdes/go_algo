package shortest_path_test

import (
	"fatdes/go_algo/shortest_path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func String(vs []interface{}) string {
	ids := make([]string, len(vs))
	for i, v := range vs {
		ids[i] = v.(*testVertex).String()
	}

	return strings.Join(ids, ",")
}

type testVertex struct {
	id    string
	edges []shortest_path.Edge
}

func (v *testVertex) Edges() []shortest_path.Edge {
	return v.edges
}

func (v *testVertex) addEdge(to *testVertex, cost int) *testVertex {
	v.edges = append(v.edges, &testEdge{
		cost: cost, from: v, to: to,
	})
	return v
}

func (v *testVertex) String() string {
	return v.id
}

type testEdge struct {
	cost int
	from *testVertex
	to   *testVertex
}

func (e *testEdge) Cost() int {
	return e.cost
}

func (e *testEdge) From() shortest_path.Vertex {
	return e.from
}

func (e *testEdge) To() shortest_path.Vertex {
	return e.to
}

type testGraph struct {
	vs map[string]*testVertex
}

func (graph *testGraph) buildTestGraph() {
	graph.vs = map[string]*testVertex{
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

func Test_UniformCost_TestEmptyFrom(t *testing.T) {
	graph := &testGraph{}
	graph.buildTestGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(nil, graph.vs["b"])
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCost_TestEmptyTo(t *testing.T) {
	graph := &testGraph{}
	graph.buildTestGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], nil)
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCost_TestNotFound(t *testing.T) {
	graph := &testGraph{}
	graph.buildTestGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], &testVertex{id: "c"})
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCost_TestSearchSameNode(t *testing.T) {
	graph := &testGraph{}
	graph.buildTestGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], graph.vs["a"])
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 0)
	assert.Equal(t, String(actual.Path), "a")
}

func Test_UniformCost_TestShortestPathFound(t *testing.T) {
	graph := &testGraph{}
	graph.buildTestGraph()

	uc := shortest_path.NewUniformCostByInterface()

	actual := uc.Find(graph.vs["a"], graph.vs["g"])
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 8)
	assert.Equal(t, String(actual.Path), "a,d,f,g")
}
