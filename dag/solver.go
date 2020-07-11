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
	solution := s.dag.Solve(problem)
	errMap := make(map[Node]error, len(problem))
	for _, nodes := range solution {
		var wg sync.WaitGroup
		var errors = make([]error, len(nodes))
		for i, node := range nodes {
			i := i
			node := node
			f := s.table[node]
			wg.Add(1)
			go func() {
				defer wg.Done()
				var err error
				// check if require failed
				for _, r := range s.dag.requires[node] {
					if errMap[r] != nil {
						err = errMap[r]
						break
					}
				}

				if err != nil {
					errors[i] = err
					return
				}
				errors[i] = f()
			}()
		}
		wg.Wait()
		// collect
		for i, err := range errors {
			errMap[nodes[i]] = err
		}
	}

	var errors = make([]error, len(problem))
	for i, node := range problem {
		errors[i] = errMap[node]
	}
	return errors
}
