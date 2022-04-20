package dag

import (
	"fmt"
	"testing"
)

const (
	FieldId Node = iota + 1
	FieldFirstName
	FieldLastName
	FieldFullName
	FieldProfile
)

type UserModel struct {
	Id        int64
	FirstName string
	LastName  string
	FullName  string
	Profile   string
}

func (m *UserModel) GetFirstName() error {
	m.FirstName = fmt.Sprintf("hello:%d", m.Id)
	return nil
}

func (m *UserModel) GetLastName() error {
	m.LastName = fmt.Sprintf("world:%d", m.Id)
	return nil
}

func (m *UserModel) GetFullName() error {
	m.FullName = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	return nil
}

func (m *UserModel) GetProfile() error {
	m.Profile = fmt.Sprintf("User:%d, with FullName: %s", m.Id, m.FullName)
	return nil
}

func (m *UserModel) Solve(n Node) error {
	switch n {
	case FieldId:
		return nil
	case FieldFirstName:
		return m.GetFirstName()
	case FieldLastName:
		return m.GetLastName()
	case FieldFullName:
		return m.GetFullName()
	case FieldProfile:
		return m.GetProfile()
	default:
		return fmt.Errorf("no such node %d", n)
	}
}

var UserModelRequires Requires = map[Node][]Node{
	FieldId:        nil,
	FieldFirstName: {FieldId},
	FieldLastName:  {FieldId},
	FieldFullName:  {FieldFirstName, FieldLastName},
	FieldProfile:   {FieldId, FieldFullName},
}

func TestSolverWithModel(t *testing.T) {
	dag, err := NewDAG(UserModelRequires)
	assertNil(t, err)
	assertNotNil(t, dag)

	user := &UserModel{Id: 1}
	solver := NewSolver(dag, user)

	fields := []Node{FieldProfile}
	err2 := solver.Solve(fields)
	assertNil(t, err2)

	assertEqual(t, "User:1, with FullName: hello:1 world:1", user.Profile)
}
