package dag

import (
	"testing"
)

func TestDAG(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0},
		2: {0, 1},
		3: {0},
		4: {2, 3},
	}

	dag, err := NewDAG(requires)
	assertNil(t, err)
	assertNotNil(t, dag)

	nodesList := dag.TopoSort()
	assertEqual(t, len(nodesList), 4)
	assertEqual(t, []Node{0}, nodesList[0])
	// assert.ElementsMatch(t, []Node{1, 3}, nodesList[1])
	assertEqual(t, []Node{2}, nodesList[2])
	assertEqual(t, []Node{4}, nodesList[3])

	problem := []Node{3}
	nodesList2 := dag.Solve((problem))
	assertEqual(t, len(nodesList2), 2)
	assertEqual(t, []Node{0}, nodesList2[0])
	assertEqual(t, []Node{3}, nodesList2[1])

	problem2 := []Node{2}
	nodesList3 := dag.Solve((problem2))
	assertEqual(t, len(nodesList3), 3)
	assertEqual(t, []Node{0}, nodesList3[0])
	assertEqual(t, []Node{1}, nodesList3[1])
	assertEqual(t, []Node{2}, nodesList3[2])
}

func TestDAGWithInvalid(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0, 2},
		2: {0, 1},
		3: {0},
		4: {2, 3},
	}
	_, err := NewDAG(requires)
	assertNotNil(t, err)
}

func TestDAGWithInvalid2(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0},
		2: {0, 1, 5},
		3: {0},
		4: {2, 3},
	}
	_, err := NewDAG(requires)
	assertNotNil(t, err)
}

func TestNodes(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0},
		2: {0, 1},
		3: {0},
		4: {2, 3},
	}
	dag, err := NewDAG(requires)
	assertNil(t, err)
	// []Node{0, 1, 3, 2, 4} or  []Node{0, 3, 1, 2, 4}
	assertEqual(t, Node(0), dag.nodes[0])
	//assert.Contains(t, []Node{1, 3}, dag.nodes[1])
	//assert.Contains(t, []Node{1, 3}, dag.nodes[2])
	assertEqual(t, Node(2), dag.nodes[3])
	assertEqual(t, Node(4), dag.nodes[4])

	assertEqual(t, dag.nodes, dag.Nodes())
}

func TestDAG3(t *testing.T) {
	dag, err := NewDAG(nil)
	assertNil(t, err)
	assertEqual(t, len(dag.Nodes()), 0)
	soloution := dag.Solve(nil)
	assertEqual(t, len(soloution), 0)
}
