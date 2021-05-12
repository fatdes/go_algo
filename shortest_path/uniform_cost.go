package shortest_path

type Result struct {
	Found bool

	Cost int
	Path []interface{}
}

type UniformCost interface {
	Find(from, to interface{}) *Result
}
