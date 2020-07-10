package dag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FieldId Node = iota +1
	FieldFirstName
	FieldLastName
	FieldFullName
	FieldProfile
)

type Model struct {
	Id int64
	FirstName string
	LastName string
	FullName string
	Profile string
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
	m.model.FullName = fmt.Sprintf("%s %s", 
		m.model.FirstName, m.model.LastName)
	return nil
}

func (m *ModelResolver) GetProfile() error {
	m.model.Profile = fmt.Sprintf("User:%d, with FullName: %s", 
		m.model.Id, m.model.FullName)
	return nil
}

func (m *ModelResolver) ResolveFactory(node Node) Call {
	r := map[Node]Call{
		FieldId: func () error {return nil},
		FieldFirstName: m.GetFirstName,
		FieldLastName: m.GetLastName,
		FieldFullName: m.GetFullName,
		FieldProfile: m.GetProfile,
	}
	return r[node]
}

func (m *ModelResolver) Deps() map[Node][]Node{
	return map[Node][]Node{
		FieldId: nil,
		FieldFirstName: {FieldId},
		FieldLastName: {FieldId},
		FieldFullName: {FieldFirstName, FieldLastName},
		FieldProfile: {FieldId, FieldFullName},
	}
}


func TestSolver(t *testing.T) {
	model := &Model{Id: 1}
	mr := &ModelResolver{model: model}
	dag, err := NewDAG(mr.Deps())
	assert.Nil(t, err)
	assert.NotNil(t, dag)

	problem := []Node{FieldProfile}
	solver := NewSolver(dag, mr.ResolveFactory)
	solver.Solve(problem)

	assert.Equal(t, "User:1, with FullName: hello:1 world:1", model.Profile)
}