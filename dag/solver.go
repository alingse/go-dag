package dag

import "sync"

type SolveFunc func() error
type SolveFuncTable map[Node]SolveFunc

type Solver struct {
	dag   *DAG
	table SolveFuncTable
}

func NewSolver(dag *DAG, table SolveFuncTable) *Solver {
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
