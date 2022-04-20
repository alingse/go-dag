package dag

import (
	"errors"
)

type DAG[Node comparable] struct {
	nodes    []Node
	requires map[Node][]Node
	topoSort [][]Node
}

var InvalidDAG = errors.New("invalid DAG")

func NewDAG[Node comparable](requires map[Node][]Node) (*DAG[Node], error) {
	ts, err := NewTopoSort(requires)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0, len(requires))
	for i := range ts {
		nodes = append(nodes, ts[i]...)
	}

	dag := &DAG[Node]{
		nodes:    nodes,
		requires: copyRequires(requires),
		topoSort: ts,
	}
	return dag, nil
}

func (d *DAG[Node]) TopoSort() [][]Node {
	return copyTopoSort(d.topoSort)
}

func (d *DAG[Node]) Nodes() []Node {
	return copyNodes(d.nodes)
}

func (d *DAG[Node]) Solve(problem []Node) [][]Node {
	needMap := make(map[Node]bool)
	for len(problem) > 0 {
		var next []Node
		for _, node := range problem {
			needMap[node] = true
			for _, r := range d.requires[node] {
				if !needMap[r] {
					next = append(next, r)
				}
			}
		}
		problem = next
	}

	var solution [][]Node
	for _, nodes := range d.topoSort {
		var rs []Node
		for _, node := range nodes {
			if needMap[node] {
				rs = append(rs, node)
			}
		}
		if len(rs) > 0 {
			solution = append(solution, rs)
		}
	}
	return solution
}

func NewTopoSort[Node comparable](requires map[Node][]Node) ([][]Node, error) {
	var ts [][]Node
	stageMap := make(map[Node]int, len(requires))
	for stage := 0; len(stageMap) < len(requires); stage++ {
		var nodes []Node
		for node, rs := range requires {
			if _, ok := stageMap[node]; ok {
				continue
			}

			var notReady []Node
			for _, r := range rs {
				if _, ok := stageMap[r]; !ok {
					notReady = append(notReady, r)
				}
			}
			if len(notReady) == 0 {
				nodes = append(nodes, node)
			}
		}
		// check got empty
		if len(nodes) == 0 {
			return nil, InvalidDAG
		}

		for _, node := range nodes {
			stageMap[node] = stage
		}

		ts = append(ts, nodes)
	}
	return ts, nil
}

func copyRequires[Node comparable](requires map[Node][]Node) map[Node][]Node {
	requires2 := make(map[Node][]Node, len(requires))
	for node, rs := range requires {
		if len(rs) > 0 {
			rs2 := make([]Node, len(rs))
			copy(rs2, rs)
			requires2[node] = rs2
		}
	}
	return requires2
}

func copyTopoSort[Node comparable](ts [][]Node) [][]Node {
	ts2 := make([][]Node, len(ts))
	for i := range ts {
		ts2[i] = make([]Node, len(ts[i]))
		copy(ts2[i], ts[i])
	}
	return ts2
}

func copyNodes[Node comparable](nodes []Node) []Node {
	nodes2 := make([]Node, len(nodes))
	copy(nodes2, nodes)
	return nodes2
}
