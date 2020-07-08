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
	assert.EqualValues(t, []Node{1, 3}, nodesList[1])
}
