package dag

import "sync"

type Call func() error
type CallFactory func(Node) Call

type Solver struct {
	dag     *DAG
	factory CallFactory
}

func NewSolver(dag *DAG, factory CallFactory) *Solver {
	return &Solver{dag: dag, factory: factory}
}

func (s *Solver) Solve(problem []Node) []error {
	errMap := make(map[Node]error, len(problem))
	solution := s.dag.Solve(problem)
	for _, nodes := range solution {
		var wg sync.WaitGroup
		for _, node := range nodes {
			wg.Add(1)
			node := node
			f := s.factory(node)
			go func() {
				defer wg.Done()
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
