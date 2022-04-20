package dag

import (
	"fmt"
)

type Solvable[Node comparable] interface {
	Solve(n Node) error
}

type Solver[Node comparable] struct {
	dag      *DAG[Node]
	solvable Solvable[Node]
}

type SolveErr[Node comparable] struct {
	Value Node
	Err   error
}

func (e SolveErr[Node]) Error() string {
	return fmt.Sprintf("solve node %v got err %s", e.Value, e.Err)
}

func NewSolver[Node comparable](dag *DAG[Node], solvable Solvable[Node]) *Solver[Node] {
	return &Solver[Node]{
		dag:      dag,
		solvable: solvable,
	}
}

func (s *Solver[Node]) Solve(problem []Node) error {
	solution := s.dag.Solve(problem)
	// fail fast
	for _, nodes := range solution {
		for _, node := range nodes {
			node := node
			err := s.solvable.Solve(node)
			if err != nil {
				return &SolveErr[Node]{Value: node, Err: err}
			}
		}
	}
	return nil
}
