package shortest_path_test

import (
	"fatdes/go_algo/shortest_path"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ByFuncString(vs []interface{}) string {
	ids := make([]string, len(vs))
	for i, v := range vs {
		ids[i] = v.(string)
	}

	return strings.Join(ids, ",")
}

type testByFuncGraph struct {
	edges     map[interface{}][]interface{}
	edgeCosts map[interface{}]int
}

func (graph *testByFuncGraph) getEdges(from interface{}) []interface{} {
	return graph.edges[from]
}

func (graph *testByFuncGraph) getEdgeEnd(edge interface{}) interface{} {
	return strings.Split(edge.(string), "_")[1]
}

func (graph *testByFuncGraph) getEdgeCost(edge interface{}) int {
	return graph.edgeCosts[edge]
}

func (graph *testByFuncGraph) addEdge(from, to interface{}, cost int) *testByFuncGraph {
	edge := fmt.Sprintf("%s_%s", from, to)
	graph.edges[from] = append(graph.edges[from], edge)
	graph.edgeCosts[edge] = cost

	return graph
}

func (graph *testByFuncGraph) buildTestByFuncGraph() {
	graph.addEdge("a", "d", 3).addEdge("a", "b", 5)
	graph.addEdge("b", "c", 1)
	graph.addEdge("c", "e", 6).addEdge("c", "g", 8)
	graph.addEdge("d", "e", 2).addEdge("d", "f", 2)
	graph.addEdge("e", "b", 4)
	graph.addEdge("f", "g", 3)
	graph.addEdge("g", "e", 4)
}

func Test_UniformCostByFunc_TestEmptyFrom(t *testing.T) {
	graph := &testByFuncGraph{edges: map[interface{}][]interface{}{}, edgeCosts: map[interface{}]int{}}
	graph.buildTestByFuncGraph()

	uc := shortest_path.NewUniformCostByFunc(graph.getEdges, graph.getEdgeEnd, graph.getEdgeCost)

	actual := uc.Find(nil, "b")
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByFunc_TestEmptyTo(t *testing.T) {
	graph := &testByFuncGraph{edges: map[interface{}][]interface{}{}, edgeCosts: map[interface{}]int{}}
	graph.buildTestByFuncGraph()

	uc := shortest_path.NewUniformCostByFunc(graph.getEdges, graph.getEdgeEnd, graph.getEdgeCost)

	actual := uc.Find("a", nil)
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByFunc_TestNotFound(t *testing.T) {
	graph := &testByFuncGraph{edges: map[interface{}][]interface{}{}, edgeCosts: map[interface{}]int{}}
	graph.buildTestByFuncGraph()

	uc := shortest_path.NewUniformCostByFunc(graph.getEdges, graph.getEdgeEnd, graph.getEdgeCost)

	actual := uc.Find("a", "h")
	assert.NotNil(t, actual)
	assert.False(t, actual.Found)
}

func Test_UniformCostByFunc_TestSearchSameNode(t *testing.T) {
	graph := &testByFuncGraph{edges: map[interface{}][]interface{}{}, edgeCosts: map[interface{}]int{}}
	graph.buildTestByFuncGraph()

	uc := shortest_path.NewUniformCostByFunc(graph.getEdges, graph.getEdgeEnd, graph.getEdgeCost)

	actual := uc.Find("a", "a")
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 0)
	assert.Equal(t, ByFuncString(actual.Path), "a")
}

func Test_UniformCostByFunc_TestShortestPathFound(t *testing.T) {
	graph := &testByFuncGraph{edges: map[interface{}][]interface{}{}, edgeCosts: map[interface{}]int{}}
	graph.buildTestByFuncGraph()

	uc := shortest_path.NewUniformCostByFunc(graph.getEdges, graph.getEdgeEnd, graph.getEdgeCost)

	actual := uc.Find("a", "g")
	assert.NotNil(t, actual)
	assert.True(t, actual.Found)
	assert.Equal(t, actual.Cost, 8)
	assert.Equal(t, ByFuncString(actual.Path), "a,d,f,g")
}
