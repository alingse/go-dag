package main

import (
	"fmt"

	"github.com/alingse/go-dag/dag"
)

type UserModel struct {
	Id        int64
	FirstName string
	LastName  string
	FullName  string
	Profile   string
}

const (
	FieldId dag.Node = iota + 1
	FieldFirstName
	FieldLastName
	FieldFullName
	FieldProfile
)

type ModelResolver struct {
	model *UserModel
}

func (m *ModelResolver) ResolveFirstName() error {
	m.model.FirstName = fmt.Sprintf("hello:%d", m.model.Id)
	return nil
}

func (m *ModelResolver) ResolveLastName() error {
	m.model.LastName = fmt.Sprintf("world:%d", m.model.Id)
	return nil
}

func (m *ModelResolver) ResolveFullName() error {
	m.model.FullName = fmt.Sprintf("%s %s", m.model.FirstName, m.model.LastName)
	return nil
}

func (m *ModelResolver) ResolveProfile() error {
	m.model.Profile = fmt.Sprintf("User:%d, with FullName: %s", m.model.Id, m.model.FullName)
	return nil
}

func (m *ModelResolver) Table() dag.SolveFuncTable {
	return map[dag.Node]dag.SolveFunc{
		FieldId:        func() error { return nil },
		FieldFirstName: m.ResolveFirstName,
		FieldLastName:  m.ResolveLastName,
		FieldFullName:  m.ResolveFullName,
		FieldProfile:   m.ResolveProfile,
	}
}

func (m *ModelResolver) Requires() map[dag.Node][]dag.Node {
	return map[dag.Node][]dag.Node{
		FieldId:        nil,
		FieldFirstName: {FieldId},
		FieldLastName:  {FieldId},
		FieldFullName:  {FieldFirstName, FieldLastName},
		FieldProfile:   {FieldId, FieldFullName},
	}
}

func main() {
	// load DAG
	var mr *ModelResolver
	d, err := dag.NewDAG(mr.Requires())
	if err != nil {
		panic(err)
	}

	// model && problem
	model := &UserModel{Id: 1}
	problem := []dag.Node{FieldProfile}

	// use mr2 Func --> set model
	mr2 := &ModelResolver{model: model}
	solver, err := dag.NewSolver(d, mr2.Table())
	if err != nil {
		panic(err)
	}

	errors := solver.Solve(problem)
	if errors[0] != nil {
		panic(errors[0])
	}

	// got the profile: 'User:1, with FullName: hello:1 world:1'
	fmt.Println(model.Profile)
}
