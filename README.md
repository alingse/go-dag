# go-dag

this repo is impelement DAG with topological sort

and add a subgraph problem Solver

## DAG

we use node dependencies to represent a DAG

build a topological-sorted nodes list based on the node dependencies.

example, the node dependencies is like this,
```
0 --> nil
1 --> 0
2 --> 0, 1
3 --> 0
4 --> 2, 3
```
and the topological sort nodes list is like this,

```
[0]
[1, 3]
[2]
[4]
```

this can solve the subgraph problem

like this.

the problem `[3]` --> got solution `[[0], [3]]`

the problem `[4]` --> got solution `[[0], [1, 3], [2], [4]]`

```
[3] --> [0] [3]
[2] --> [0] [1] [2]
[1, 3] --> [0] [1, 3]
[2, 3] --> [0] [1, 3] [2]
[2, 4] --> [0] [1, 3] [2] [4]
```

## Solver

`DAG(nodes) ----> TopoSort [][]Node --> problem: []Node --> solution [][]Node `

`Solution` is the projection (subgraph) of the TopoSort result under the Problem


## Solver

impelement `Solvable`

```go
type Solvable interface {
	Solve(n Node) error
}
```

### example

```go
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

var fieldDeps = map[string][]string{
	FieldId:        nil,
	FieldFirstName: {FieldId},
	FieldLastName:  {FieldId},
	FieldFullName:  {FieldFirstName, FieldLastName},
	FieldProfile:   {FieldId, FieldFullName},
}

func main() {
	userDAG, err := dag.NewDAG(fieldDeps)
	if err != nil {
		panic(err)
	}

	user := &UserModel{Id: 1}
	userSolver := dag.NewSolver[string](userDAG, user)

	fields := []string{FieldProfile}
	_ = userSolver.Solve(fields)
	// I'm User 1, my FullName is「firstName1 lastName1」
	fmt.Println(user.Profile)
}
```

so, `userSolver` can auto solve the `UserModel` with the input `fields`

```go
user := &UserModel{Id: 1}
userSolver := dag.NewSolver[string](userDAG, user)

fields := []string{FieldProfile}
_ = userSolver.Solve(fields)
```

# TODO

目前的 Solve 是 failfast

也许每一层里面可以尽力执行, 或者将 Solveable 作为多 node 的 `Solve(nodes []Node) error`

这样具体是并发、fail fast、try all 都由使用者决定
