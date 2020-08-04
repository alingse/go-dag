package dag

import (
	"errors"
)

type Node int

type DAG struct {
	nodes    []Node
	requires map[Node][]Node
	topoSort [][]Node
}

var InvalidDAG = errors.New("invalid DAG")

func NewDAG(requires map[Node][]Node) (*DAG, error) {
	requires2 := make(map[Node][]Node, len(requires))
	for node, rs := range requires {
		nodes := make([]Node, len(rs))
		copy(nodes, rs)
		requires2[node] = nodes
	}
	requires = requires2 // copy

	ts, err := topoSort(requires)
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0, len(requires))
	for i := range ts {
		nodes = append(nodes, ts[i]...)
	}

	dag := &DAG{
		nodes:    nodes,
		requires: requires,
		topoSort: ts,
	}
	return dag, nil
}

func topoSort(requires map[Node][]Node) ([][]Node, error) {
	// check all nodes has required
	for _, rs := range requires {
		for _, r := range rs {
			if _, ok := requires[r]; !ok {
				return nil, InvalidDAG
			}
		}
	}

	var ts [][]Node
	stageMap := make(map[Node]int, len(requires))
	for stage := 0; len(stageMap) < len(requires); stage++ {
		var nodes []Node
		for node, rs := range requires {
			if _, ok := stageMap[node]; ok {
				continue
			}

			var frs []Node
			for _, r := range rs {
				if _, ok := stageMap[r]; !ok {
					frs = append(frs, r)
				}
			}
			if len(frs) == 0 {
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

func (d *DAG) TopoSort() [][]Node {
	ts := make([][]Node, len(d.topoSort))
	for i := range d.topoSort {
		ts[i] = make([]Node, len(d.topoSort[i]))
		copy(ts[i], d.topoSort[i])
	}
	return ts
}

func (d *DAG) Nodes() []Node {
	nodes := make([]Node, len(d.nodes))
	copy(nodes, d.nodes)
	return nodes
}

func (d *DAG) Solve(problem []Node) [][]Node {
	need := make(map[Node]bool)
	for len(problem) > 0 {
		var next []Node
		for _, node := range problem {
			need[node] = true
			for _, r := range d.requires[node] {
				if !need[r] {
					next = append(next, r)
				}
			}
		}
		problem = next
	}

	var soloution [][]Node
	for _, nodes := range d.topoSort {
		var rs []Node
		for _, node := range nodes {
			if need[node] {
				rs = append(rs, node)
			}
		}
		if len(rs) > 0 {
			soloution = append(soloution, rs)
		}
	}
	return soloution
}
