package main

import (
	"fmt"

	"github.com/alingse/go-dag"
)

const (
	FieldId        = "id"
	FieldFirstName = "first_name"
	FieldLastName  = "last_name"
	FieldFullName  = "full_name"
	FieldProfile   = "profile"
)

type UserModel struct {
	Id        int64
	FirstName string
	LastName  string
	FullName  string
	Profile   string
}

func (m *UserModel) GetFirstName() error {
	m.FirstName = fmt.Sprintf("firstName%d", m.Id)
	return nil
}

func (m *UserModel) GetLastName() error {
	m.LastName = fmt.Sprintf("lastName%d", m.Id)
	return nil
}

func (m *UserModel) GetFullName() error {
	m.FullName = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	return nil
}

func (m *UserModel) GetProfile() error {
	m.Profile = fmt.Sprintf("I'm User %d, my FullName is「%s」", m.Id, m.FullName)
	return nil
}

func (m *UserModel) Solve(n string) error {
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
		return fmt.Errorf("no such node %v", n)
	}
}

var filedRequires = map[string][]string{
	FieldId:        nil,
	FieldFirstName: {FieldId},
	FieldLastName:  {FieldId},
	FieldFullName:  {FieldFirstName, FieldLastName},
	FieldProfile:   {FieldId, FieldFullName},
}

func main() {
	userDAG, err := dag.NewDAG(filedRequires)
	if err != nil {
		panic(err)
	}

	user := &UserModel{Id: 1}
	userSolver := dag.NewSolver[string](userDAG, user)

	fields := []string{FieldProfile}
	err = userSolver.Solve(fields)
	if err != nil {
		panic(err)
	}

	// I'm User 1, my FullName is「firstName1 lastName1」
	fmt.Println(user.Profile)
}
