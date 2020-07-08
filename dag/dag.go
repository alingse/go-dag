package dag

import (
	"errors"
)

type Node int

type DAG struct {
	requires map[Node][]Node
	topSort  [][]Node
}

func NewDAG(requires map[Node][]Node) (*DAG, error) {
	ts, err := topSort(requires)
	if err != nil {
		return nil, err
	}
	return &DAG{requires: requires, topSort: ts}, nil
}

var InvalidDAG = errors.New("invalid DAG")

func topSort(requires map[Node][]Node) ([][]Node, error) {
	if len(requires) == 0 {
		return nil, nil
	}
	var stageNodes [][]Node
	stageMap := make(map[Node]int, len(requires))
	for stage := 0; len(stageMap) < len(requires); stage++ {
		var nodes []Node
		for node, require := range requires {
			// checked
			if _, ok := stageMap[node]; ok {
				continue
			}
			// fitler
			var rs []Node
			for _, r := range require {
				if _, ok := stageMap[r]; !ok {
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
		// this stage
		for _, node := range nodes {
			stageMap[node] = stage
		}
		stageNodes = append(stageNodes, nodes)
	}
	return stageNodes, nil
}

func (d *DAG) TopSort() [][]Node {
	return d.topSort
}
