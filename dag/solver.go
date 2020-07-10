package dag

import "sync"

type Call func() error
type CallTable map[Node]Call

type Solver struct {
	dag   *DAG
	table CallTable
}

func NewSolver(dag *DAG, table CallTable) *Solver {
	return &Solver{dag: dag, table: table}
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
