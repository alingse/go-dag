package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	assert.NotNil(t, dag)

	nodesList := dag.TopoSort()
	assert.Len(t, nodesList, 4)
	assert.Equal(t, []Node{0}, nodesList[0])
	assert.ElementsMatch(t, []Node{1, 3}, nodesList[1])
	assert.Equal(t, []Node{2}, nodesList[2])
	assert.Equal(t, []Node{4}, nodesList[3])

	problem := []Node{3}
	nodesList2 := dag.Solve((problem))
	assert.Len(t, nodesList2, 2)
	assert.Equal(t, []Node{0}, nodesList2[0])
	assert.Equal(t, []Node{3}, nodesList2[1])

	problem2 := []Node{2}
	nodesList3 := dag.Solve((problem2))
	assert.Len(t, nodesList3, 3)
	assert.Equal(t, []Node{0}, nodesList3[0])
	assert.Equal(t, []Node{1}, nodesList3[1])
	assert.Equal(t, []Node{2}, nodesList3[2])
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
	assert.NotNil(t, err)
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
	assert.NotNil(t, err)
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
	assert.Nil(t, err)
	// []Node{0, 1, 3, 2, 4} or  []Node{0, 3, 1, 2, 4}
	assert.Equal(t, Node(0), dag.nodes[0])
	assert.Contains(t, []Node{1, 3}, dag.nodes[1])
	assert.Contains(t, []Node{1, 3}, dag.nodes[2])
	assert.Equal(t, Node(2), dag.nodes[3])
	assert.Equal(t, Node(4), dag.nodes[4])

	assert.Equal(t, dag.nodes, dag.Nodes())
}

func TestDAG3(t *testing.T) {
	dag, err := NewDAG(nil)
	assert.Nil(t, err)
	assert.Len(t, dag.Nodes(), 0)
	soloution := dag.Solve(nil)
	assert.Len(t, soloution, 0)
}
