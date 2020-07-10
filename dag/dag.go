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
	ts, err := topoSort(requires)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0, len(requires))
	for node := range requires {
		nodes = append(nodes, node)
	}
	dag := &DAG{
		nodes:    nodes,
		requires: requires,
		topoSort: ts,
	}
	return dag, nil
}

func topoSort(requires map[Node][]Node) ([][]Node, error) {
	if len(requires) == 0 {
		return nil, nil
	}
	// check all nodes has required
	for _, rs := range requires {
		for _, r := range rs {
			if _, ok := requires[r]; !ok {
				return nil, InvalidDAG
			}
		}
	}

	var nodesList [][]Node
	nodeMap := make(map[Node]int, len(requires))
	for stage := 0; len(nodeMap) < len(requires); stage++ {
		var nodes []Node
		for node, require := range requires {
			if _, ok := nodeMap[node]; ok {
				continue
			}
			// fitler
			var rs []Node
			for _, r := range require {
				if _, ok := nodeMap[r]; !ok {
					rs = append(rs, r)
				}
			}
			if len(rs) == 0 {
				nodes = append(nodes, node)
			}
		}
		// check got empty
		if len(nodes) == 0 {
			return nil, InvalidDAG
		}

		for _, node := range nodes {
			nodeMap[node] = stage
		}
		nodesList = append(nodesList, nodes)
	}
	return nodesList, nil
}

func (d *DAG) TopoSort() [][]Node {
	// return a copy of d.topoSort
	var ts = make([][]Node, len(d.topoSort))
	for i := range d.topoSort {
		ts[i] = make([]Node, len(d.topoSort[i]))
		copy(ts[i], d.topoSort[i])
	}
	return ts
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

	var solve [][]Node
	for _, nodes := range d.topoSort {
		var rs []Node
		for _, node := range nodes {
			if need[node] {
				rs = append(rs, node)
			}
		}
		if len(rs) > 0 {
			solve = append(solve, rs)
		}
	}
	return solve
}

func (d *DAG) Nodes() []Node {
	return d.nodes
}
