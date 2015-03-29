package lp

import "fmt"

// Carries an ordinary primal dictionary to solution.
// Assumes that initial dictionary is feasible.
func Solve(dict *Dict) (final *Dict, unbnd bool) {
	return SolveEps(dict, DefaultEps)
}

func SolveEps(dict *Dict, eps float64) (final *Dict, unbnd bool) {
	if !dict.Feas() {
		panic("initial dictionary infeasible")
	}

	// Pivot until reaching the solution.
	var iter int
	for {
		piv := NextBlandEps(dict, eps)
		if piv.Unbounded {
			return dict, true
		}
		if piv.Final {
			return dict, false
		}
		dict = dict.Pivot(piv.Enter, piv.Leave)
		iter++
		fmt.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}
}

// Solves the feasibility problem of a given dictionary.
// If feasible, returns a feasible dictionary of the original problem.
// Assumes that the original dictionary is infeasible.
func SolveFeas(dict *Dict) (final *Dict, infeas bool) {
	return SolveFeasEps(dict, DefaultEps)
}

func SolveFeasEps(orig *Dict, eps float64) (feas *Dict, infeas bool) {
	// Transform to a dictionary for the feasibility problem.
	dict := ToFeasDict(orig)

	// Perform feasibility pivots.
	var iter int
	for {
		piv := NextFeasBlandEps(dict, eps)
		if piv.Unbounded {
			// Auxiliary problem
			//   min  x  s.t.  x >= 0, ...
			// is always bounded.
			panic("unbounded")
		}
		if piv.Final {
			break
		}
		dict = dict.Pivot(piv.Enter, piv.Leave)
		iter++
		fmt.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}

	// The gap to feasibility such that (A x - u 1 <= b).
	u := -dict.Obj()
	if u > eps {
		return nil, true
	}

	// Transform back to a feasible dictionary for the original problem.
	return FromFeasDict(dict, orig), false
}