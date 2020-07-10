# go-dag
try implement DAG resolver and topological sort

## DAG

根据依赖关系构建拓扑排序

```
0 --> nil
1 --> 0
2 --> 0, 1
3 --> 0
4 --> 2, 3
```

得到的拓扑排序就是
[0] [1, 3] [2] [4]

### Solve

将传入的 []Node 作为问题 Solve 一种解决方案

```
[3] --> [0] [3]
[2] --> [0] [1] [2]
[1, 3] --> [0] [1, 3]
[2, 3] --> [0] [1, 3] [2]
[2, 4] --> [0] [1, 3] [2] [4]
```

solve 数组决定依赖解决顺序, solve 第 i 项可并发处理

## Solver

根据构建好的 DAG 和 SolveFuncTable 去执行 problem 的求解过程

如 examples 给出的例子

```go
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

func (m *ModelResolver) Requires() dag.DAGRequires {
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
```

依次声明 `Model` 和 `Field` 以及各个字段的具体 `ResolveXXX` 实现

ModelResolver 给出 `Requires()` 和 `Table()`

给出具体要处理的 problem --> `[]{FieldProfile}`

solver 将自动 Solve

### Solve

目前 Solver.Solve 函数是尽力执行完所有可能执行的。

也许可以提供碰见任意 err 就停止的 Solve (failfast?), 不过这种情况 err 通过依赖路径一路 pop, 最终返回 problem 对应的 err 优点的难度
