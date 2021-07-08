package dag

import (
	"errors"
	"sync"
)

type (
	SolveFunc func() error
)

type Solvable interface {
	GetSolveFunc(n Node) SolveFunc
}

type Solver struct {
	dag      *DAG
	solvable Solvable
}

var InvalidSolver = errors.New("invalid Solver")

func NewSolver(dag *DAG, solvable Solvable) (*Solver, error) {
	solver := &Solver{
		dag:      dag,
		solvable: solvable,
	}
	return solver, nil
}

func (s *Solver) Solve(problem []Node) []error {
	solution := s.dag.Solve(problem)
	errMap := make(map[Node]error, len(problem))
	for _, nodes := range solution {
		var wg sync.WaitGroup
		errors := make([]error, len(nodes))
		for i, node := range nodes {
			i := i
			node := node
			f := s.solvable.GetSolveFunc(node)
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

	errors := make([]error, len(problem))
	for i, node := range problem {
		errors[i] = errMap[node]
	}
	return errors
}
