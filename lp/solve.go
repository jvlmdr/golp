package lp

import "fmt"

// Solve solves a linear program.
func Solve(dict *Dict) (final *Dict, err error) {
	return SolveEps(dict, DefaultEps)
}

func SolveEps(dict *Dict, eps float64) (final *Dict, err error) {
	if !dict.Feas() {
		// If the solution associated with the dictionary is infeasible,
		// attempt find a feasible dictionary.
		var infeas bool
		dict, infeas = SolveFeasEps(dict, eps)
		if infeas {
			return nil, fmt.Errorf("infeasible problem")
		}
	}
	var unbnd bool
	dict, unbnd = PivotToFinalEps(dict, eps)
	if unbnd {
		return nil, fmt.Errorf("unbounded problem")
	}
	return dict, nil
}

// PivotToFinal carries a feasible dictionary to solution.
// Assumes that initial dictionary is feasible.
func PivotToFinal(dict *Dict) (final *Dict, unbnd bool) {
	return PivotToFinalEps(dict, DefaultEps)
}

func PivotToFinalEps(dict *Dict, eps float64) (final *Dict, unbnd bool) {
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
		//log.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}
}

// SolveFeas solves the feasibility problem of a given dictionary.
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
		//log.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}

	// The gap to feasibility such that (A x - u 1 <= b).
	u := -dict.Obj()
	if u > eps {
		return nil, true
	}

	// Transform back to a feasible dictionary for the original problem.
	return FromFeasDict(dict, orig), false
}
