# go-dag

go DAG with topological sort, and subgraph problem Solver (with some child nodes)

## DAG

we use node dependencies to represent a DAG

and build a topological-sorted nodes based on the node dependencies.

example, the node dependencies is like this,
```
0 --> nil
1 --> 0
2 --> 0, 1
3 --> 0
4 --> 2, 3
```
and the topological sort nodes like this,

```
[0]
[1, 3]
[2]
[4]
```

### Solve

and this DAG can solve the subgraph problem

like

the problem `[3]` --> got solution `[[0], [3]]`

the problem `[4]` --> got solution `[[0], [1, 3], [2], [4]]`

```
[3] --> [0] [3]
[2] --> [0] [1] [2]
[1, 3] --> [0] [1, 3]
[2, 3] --> [0] [1, 3] [2]
[2, 4] --> [0] [1, 3] [2] [4]
```

有时, 不需要全部解决 DAG 的每个 Node

给定 []None 作为一个 Problem, 对拓扑排序的结果 TopoSort 进行抽取

得到 TopoSort 结果在该 Problem 下的投影, 就是本次需要的一个解决方案 (Solution)。

比如


Problem [3] 的解决方案是 [0], [3] 只需要依次处理 0, 3 节点即可

Problem [2, 4] 的解决方案是 [0] [1, 3] [2] [4]

solution 的数组决定解决顺序, 每一层内部可以并行处理。

## Solver

定义 `Solvable` 作为一个实际可解决的问题

```go
type Solvable interface {
	Solve(n Node) error
}
```

根据构建好的 DAG 和 Solvable 去执行 problem 的求解过程

如 examples 给出的例子

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
```

依次声明 `UserModel` 的各个字段 `Field` 的 Resolve Func, 并实现 Solvable 接口 `Solve`

根据声明的 UserDAG 和 UserModel 构建 solver 加上传入 fields[FieldProfile]
最后就能自动 Solve 得到 UserModel.Profile 字段


```go
user := &UserModel{Id: 1}
userSolver := dag.NewSolver(userDAG, user)

fields := []dag.Node{FieldProfile}
err = userSolver.Solve(fields)
```

### SolveType

目前的 Solve 是 failfast

也许每一层里面可以尽力执行, 或者将 Solveable 作为多 node 的 `Solve(nodes []Node) error`

这样具体是并发、fail fast、try all 都由使用者决定
