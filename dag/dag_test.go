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

	nodesList := dag.TopSort()
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
