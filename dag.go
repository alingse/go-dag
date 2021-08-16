package dag

import (
	"errors"
)

type Node int
type TopoSort [][]Node
type Requires map[Node][]Node

type DAG struct {
	nodes    []Node
	requires Requires
	topoSort TopoSort
}

var InvalidDAG = errors.New("invalid DAG")

func NewDAG(requires Requires) (*DAG, error) {
	ts, err := NewTopoSort(requires)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0, len(requires))
	for i := range ts {
		nodes = append(nodes, ts[i]...)
	}

	dag := &DAG{
		nodes:    nodes,
		requires: copyRequires(requires),
		topoSort: ts,
	}
	return dag, nil
}

func (d *DAG) TopoSort() [][]Node {
	return copyTopoSort(d.topoSort)
}

func (d *DAG) Nodes() []Node {
	return copyNodes(d.nodes)
}

func (d *DAG) Solve(problem []Node) TopoSort {
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

	var solution TopoSort
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

func NewTopoSort(requires Requires) (TopoSort, error) {
	var ts TopoSort
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

func copyRequires(requires Requires) Requires {
	requires2 := make(Requires, len(requires))
	for node, rs := range requires {
		if len(rs) > 0 {
			rs2 := make([]Node, len(rs))
			copy(rs2, rs)
			requires2[node] = rs2
		}
	}
	return requires2
}

func copyTopoSort(ts TopoSort) TopoSort {
	ts2 := make(TopoSort, len(ts))
	for i := range ts {
		ts2[i] = make([]Node, len(ts[i]))
		copy(ts2[i], ts[i])
	}
	return ts2
}

func copyNodes(nodes []Node) []Node {
	nodes2 := make([]Node, len(nodes))
	copy(nodes2, nodes)
	return nodes2
}
