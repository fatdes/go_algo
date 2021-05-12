package shortest_path

type Result struct {
	Found bool

	Cost int
	Path []interface{}
}

type UniformCost interface {
	Find(from, to interface{}) *Result
}

type node struct {
	vertex    interface{}
	totalCost int
	path      []interface{}
}

func (n *node) cost() int {
	return n.totalCost
}
