# go-dag

implement DAG with topological sort

and DAG Solve help code

## DAG

根据依赖关系构建拓扑排序

整个 DAG 依赖按照拓扑排序依次从前到后递进

例子

```
0 --> nil
1 --> 0
2 --> 0, 1
3 --> 0
4 --> 2, 3
```

得到的拓扑排序就是
[0] [1, 3] [2] [4]

需要先有 [0] 然后才可以得到 [1, 3], 再得到 [2] 最后才能得到 [4]

一层中的节点没有依赖关系, 下一层只能等到上一层完全就绪才可以。

### Solve

有时, 不需要全部解决 DAG 的每个 Node

给定 []None 作为一个 Problem, 对拓扑排序的结果 TopoSort 进行抽取

得到 TopoSort 结果在该 Problem 下的投影, 就是本次需要的一个解决方案 (Solution)。

比如

```
[3] --> [0] [3]
[2] --> [0] [1] [2]
[1, 3] --> [0] [1, 3]
[2, 3] --> [0] [1, 3] [2]
[2, 4] --> [0] [1, 3] [2] [4]
```

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
	FieldId dag.Node = iota + 1
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

func (m *UserModel) Solve(n dag.Node) error {
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

var UserModelRequires dag.Requires = map[dag.Node][]dag.Node{
	FieldId:        nil,
	FieldFirstName: {FieldId},
	FieldLastName:  {FieldId},
	FieldFullName:  {FieldFirstName, FieldLastName},
	FieldProfile:   {FieldId, FieldFullName},
}

func main() {
	userDAG, err := dag.NewDAG(UserModelRequires)
	if err != nil {
		panic(err)
	}

	user := &UserModel{Id: 1}
	userSolver := dag.NewSolver(userDAG, user)

	// fields
	fields := []dag.Node{FieldProfile}
	err = userSolver.Solve(fields)
	if err != nil {
		panic(err)
	}

	// got the profile: 'User:1, with FullName: hello:1 world:1'
	fmt.Println(user.Profile)
}
```

依次声明 `UserModel` 的各个字段 `Field` 的 Func, 并实现 Solvable 接口


```go
user := &UserModel{Id: 1}
userSolver := dag.NewSolver(userDAG, user)

fields := []dag.Node{FieldProfile}
err = userSolver.Solve(fields)
```

根据声明的 UserDAG 和 UserModel 构建 solver 加上传入 fields[FieldProfile]

最后就能自动 Solve 得到 UserModel.Profile 字段

### SolveType

目前的 Solve 是 failfast

也许每一层里面可以尽力执行, 或者将 Solveable 作为多 node 的 `Solve(nodes []Node) error`

这样具体是并发、fail fast、try all 都有使用者决定
