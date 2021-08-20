package dag

import (
	"fmt"
)

type Solvable interface {
	Solve(n Node) error
}

type Solver struct {
	dag      *DAG
	solvable Solvable
}

type SolveErr struct {
	Node
	Err error
}

func (e SolveErr) Error() string {
	return fmt.Sprintf("solve node %d got err %s", e.Node, e.Err)
}

func NewSolver(dag *DAG, solvable Solvable) *Solver {
	return &Solver{
		dag:      dag,
		solvable: solvable,
	}
}

func (s *Solver) Solve(problem []Node) error {
	solution := s.dag.Solve(problem)
	// fail fast
	for _, nodes := range solution {
		for _, node := range nodes {
			node := node
			err := s.solvable.Solve(node)
			if err != nil {
				return &SolveErr{Node: node, Err: err}
			}
		}
	}
	return nil
}
