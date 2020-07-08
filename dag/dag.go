package dag

type Node int

type DAG struct {
	requires map[Node][]Node
}

func NewDAG(requires map[Node][]Node) (*DAG, error) {
	// TODO: check is ok
	return &DAG{requires: requires}, nil
}

func (d *DAG) TopSort() [][]Node {
	if len(d.requires) == 0 {
		return nil
	}

	stageNodes := make([][]Node, 0, 1)
	stageMap := make(map[Node]int, len(d.requires))
	for stage := 0; len(stageMap) < len(d.requires); stage++ {
		var nodes []Node
		for node, requires := range d.requires {
			// checked
			if _, ok := stageMap[node]; ok {
				continue
			}
			// fitler
			var rs []Node
			for _, r := range requires {
				if _, ok := stageMap[r]; !ok {
					rs = append(rs, r)
				}
			}
			if len(rs) == 0 {
				nodes = append(nodes, node)
			}
		}
		// this stage
		for _, node := range nodes {
			stageMap[node] = stage
		}
		stageNodes = append(stageNodes, nodes)
	}
	return stageNodes
}
