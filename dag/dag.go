package dag

import (
	"errors"
)

type Node int
type DAGRequires map[Node][]Node

type DAG struct {
	nodes    []Node
	requires DAGRequires
	topoSort [][]Node
}

var InvalidDAG = errors.New("invalid DAG")

func NewDAG(requires DAGRequires) (*DAG, error) {
	// copy
	requires2 := make(DAGRequires, len(requires))
	for node, rs := range requires {
		nodes := make([]Node, len(rs))
		copy(nodes, rs)
		requires2[node] = nodes
	}
	requires = requires2

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

func topoSort(requires DAGRequires) ([][]Node, error) {
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

func (d *DAG) Nodes() []Node {
	// return a copy of d.nodes
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
