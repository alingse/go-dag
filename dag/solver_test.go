package dag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FieldId Node = iota + 1
	FieldFirstName
	FieldLastName
	FieldFullName
	FieldProfile
)

type Model struct {
	Id        int64
	FirstName string
	LastName  string
	FullName  string
	Profile   string
}

type ModelResolver struct {
	model *Model
}

func (m *ModelResolver) GetFirstName() error {
	m.model.FirstName = fmt.Sprintf("hello:%d", m.model.Id)
	return nil
}

func (m *ModelResolver) GetLastName() error {
	m.model.LastName = fmt.Sprintf("world:%d", m.model.Id)
	return nil
}

func (m *ModelResolver) GetFullName() error {
	m.model.FullName = fmt.Sprintf("%s %s", m.model.FirstName, m.model.LastName)
	return nil
}

func (m *ModelResolver) GetProfile() error {
	m.model.Profile = fmt.Sprintf("User:%d, with FullName: %s", m.model.Id, m.model.FullName)
	return nil
}

func (m *ModelResolver) ResolveTable() map[Node]SolveFunc {
	return map[Node]SolveFunc{
		FieldId:        func() error { return nil },
		FieldFirstName: m.GetFirstName,
		FieldLastName:  m.GetLastName,
		FieldFullName:  m.GetFullName,
		FieldProfile:   m.GetProfile,
	}
}

func (m *ModelResolver) GetSolveFunc(node Node) SolveFunc {
	table := m.ResolveTable()
	return table[node]
}

func (m *ModelResolver) ResolveDeps() map[Node][]Node {
	return map[Node][]Node{
		FieldId:        nil,
		FieldFirstName: {FieldId},
		FieldLastName:  {FieldId},
		FieldFullName:  {FieldFirstName, FieldLastName},
		FieldProfile:   {FieldId, FieldFullName},
	}
}

func TestSolverWithModel(t *testing.T) {
	model := &Model{Id: 1}
	mr := &ModelResolver{model: model}
	dag, err := NewDAG(mr.ResolveDeps())
	assert.Nil(t, err)
	assert.NotNil(t, dag)

	problem := []Node{FieldProfile}
	solver, err := NewSolver(dag, mr)
	assert.Nil(t, err)
	assert.NotNil(t, solver)

	solver.Solve(problem)
	assert.Equal(t, "User:1, with FullName: hello:1 world:1", model.Profile)
}

func TestSolve1(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0},
		2: {0, 1},
		3: {0, 1},
		4: {3, 1},
	}
	dag, err := NewDAG(requires)
	assert.Nil(t, err)

	var solvable Solvable
	_, err = NewSolver(dag, solvable)
	assert.Nil(t, err)
}

type SolveTable map[Node]SolveFunc

func (x SolveTable) GetSolveFunc(node Node) SolveFunc {
	return x[node]
}

func TestSolve2(t *testing.T) {
	requires := map[Node][]Node{
		0: nil,
		1: {0},
		2: {0, 1},
		3: {0, 1},
		4: {3, 1},
	}
	dag, err := NewDAG(requires)
	assert.Nil(t, err)

	table := SolveTable{
		0: func() error { return nil },
		1: func() error { return nil },
		2: func() error { return nil },
		3: func() error { return fmt.Errorf("panic") },
		4: func() error { return nil },
	}

	solver, err := NewSolver(dag, table)
	assert.Nil(t, err)

	errors := solver.Solve([]Node{2})
	assert.Nil(t, errors[0])

	errors2 := solver.Solve([]Node{2, 4})
	assert.Nil(t, errors2[0])
	assert.NotNil(t, errors2[1])
}
