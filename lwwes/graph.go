package lwwes

import (
	"fmt"
)

// Node with id, same id means same node
type Node struct {
	ID string
}

var nodeEqual = func(d1 interface{}, d2 interface{}) bool {
	return d1.(*Node).ID == d2.(*Node).ID
}

func (node *Node) String() string {
	return fmt.Sprintf("node: %s", node.ID)
}

// Edge between two nodes, direction from -> to
type Edge struct {
	From *Node
	To   *Node
}

var edgeEqual = func(d1 interface{}, d2 interface{}) bool {
	e1 := d1.(*Edge)
	e2 := d2.(*Edge)

	return nodeEqual(e1.From, e2.From) && nodeEqual(e1.To, e2.To)
}

func (edge *Edge) String() string {
	return fmt.Sprintf("edge: (%s) -> (%s)", edge.From, edge.To)
}

// TreeNode of Node -> Children
type TreeNode struct {
	Node     *Node
	Children []*TreeNode
}

func (tn *TreeNode) String() string {
	return fmt.Sprintf("(%s) -> [%s]", tn.Node, tn.Children)
}

func depthFirstSearch(current *TreeNode, f func(current *TreeNode) bool) *TreeNode {
	if f(current) {
		return current
	}

	for _, child := range current.Children {
		if found := depthFirstSearch(child, f); found != nil {
			return found
		}
	}

	return nil
}

// RootNode is always there
var RootNode = &Node{ID: "__root"}

// Graph with static root node and node/edge LWW-Element-Sets
type Graph struct {
	Root *Node

	Nodes *ElementSet
	Edges *ElementSet

	dataEqual Equal
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	graph := &Graph{
		Nodes:     NewElementSet(NewSet(nodeEqual), NewSet(nodeEqual)),
		Edges:     NewElementSet(NewSet(edgeEqual), NewSet(edgeEqual)),
		dataEqual: nodeEqual,
	}
	graph.Nodes.AddData(RootNode, 1)
	return graph
}

func contains(a []string, s string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == s {
			return true
		}
	}

	return false
}

func (g *Graph) buildRootedTree() *TreeNode {
	// orphan node if there is no edge linked to it OR node is NOT in the Node Set
	orphanNodeIDs := []string{}
	updateOrphanNode := func(node *Node) bool {
		if contains(orphanNodeIDs, node.ID) {
			return true
		}

		if g.Nodes.Lookup(node) == nil {
			orphanNodeIDs = append(orphanNodeIDs, node.ID)
			return true
		}

		return false
	}

	// count how many links to a node
	childrenCounts := map[string]int{}
	updateChildrenCount := func(node *Node) int {
		if _, found := childrenCounts[node.ID]; !found {
			childrenCounts[node.ID] = 1
		} else {
			childrenCounts[node.ID]++
		}
		return childrenCounts[node.ID]
	}

	edges := g.Edges.LookupAll()
	edgeMap := map[string][]*Edge{}
	for _, e := range edges {
		edge := e.Data.(*Edge)
		updateOrphanNode(edge.From)
		updateOrphanNode(edge.To)
		updateChildrenCount(edge.To)
	}

	for _, e := range edges {
		edge := e.Data.(*Edge)
		if
		// connection policy: skip
		contains(orphanNodeIDs, edge.From.ID) || contains(orphanNodeIDs, edge.To.ID) ||
			// mapping policy: zero
			childrenCounts[edge.From.ID] > 1 || childrenCounts[edge.To.ID] > 1 {
			continue
		}

		edgeMap[edge.From.ID] = append(edgeMap[edge.From.ID], edge)
	}

	var buildTree func(node *Node) *TreeNode

	buildTree = func(node *Node) *TreeNode {
		children := edgeMap[node.ID]

		treeNode := &TreeNode{
			Node:     node,
			Children: make([]*TreeNode, len(children)),
		}

		for i, child := range children {
			treeNode.Children[i] = buildTree(child.To)
		}

		return treeNode
	}

	return buildTree(RootNode)
}

// AddNode to the graph
func (g *Graph) AddNode(node *Node, parent *Node, timestamp int64) {
	g.Nodes.AddData(node, timestamp)
	g.Edges.AddData(&Edge{From: parent, To: node}, timestamp)
}

// RemoveNode from the graph
func (g *Graph) RemoveNode(node *Node, parent *Node, timestamp int64) {
	g.Nodes.RemoveData(node, timestamp)
	g.Edges.RemoveData(&Edge{From: parent, To: node}, timestamp)
}

// QueryNode of id
func (g *Graph) QueryNode(id string) *TreeNode {
	// shortcut
	nodes := g.Nodes.LookupAll()
	found := false
	for _, node := range nodes {
		if node.Data.(*Node).ID == id {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	root := g.buildRootedTree()
	node := depthFirstSearch(root, func(current *TreeNode) bool {
		return current.Node.ID == id
	})

	return node
}

// QueryAllNodes get all the available nodes
// connection policy: skip orphan nodes
// mapping policy: not shown if multiple parents
func (g *Graph) QueryAllNodes() *TreeNode {
	return g.buildRootedTree()
}

// Merge this graph with other graph
func (g *Graph) Merge(other *Graph) {
	g.Nodes.Merge(other.Nodes)
	g.Edges.Merge(other.Edges)
}
