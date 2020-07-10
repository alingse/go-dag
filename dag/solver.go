package dag

import (
	"errors"
	"sync"
)

type SolveFunc func() error
type SolveFuncTable map[Node]SolveFunc

type Solver struct {
	dag   *DAG
	table SolveFuncTable
}

var InvalidSolver = errors.New("invalid Solver")

func NewSolver(dag *DAG, table SolveFuncTable) (*Solver, error) {
	for node := range dag.requires {
		if table[node] == nil {
			return nil, InvalidSolver
		}
	}
	solver := &Solver{dag: dag, table: table}
	return solver, nil
}

func (s *Solver) Solve(problem []Node) []error {
	errMap := make(map[Node]error, len(problem))
	solution := s.dag.Solve(problem)
	for _, nodes := range solution {
		var wg sync.WaitGroup
		for _, node := range nodes {
			wg.Add(1)
			node := node
			f := s.table[node]
			go func() {
				defer wg.Done()
				// require failed
				for _, r := range s.dag.requires[node] {
					if errMap[r] != nil {
						errMap[node] = errMap[r]
						return
					}
				}
				errMap[node] = f()
			}()
		}
		wg.Wait()
	}

	var errs []error
	for _, node := range problem {
		errs = append(errs, errMap[node])
	}
	return errs
}
