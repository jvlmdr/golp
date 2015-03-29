package lp

// Finds index (not label) of next non-basic variable to enter
// under Bland's rule.
func toEnterBland(dict *Dict, eps float64) (enter int, final bool) {
	// Find variable with lowest index.
	var (
		found bool
		arg   int
		min   int
	)

	// Find lowest-index variable with positive objective coefficient.
	for i := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[i] <= eps {
			continue
		}

		index := dict.NonBasic[i]
		if !found || index < min {
			found = true
			arg, min = i, index
		}
	}
	return arg, !found
}

// Finds index (not label) of next basic variable to leave
// under Bland's rule given non-basic variable to enter.
//
// Assumes that dictionary is feasible.
func toLeaveBland(dict *Dict, enter int, eps float64) (leave int, unbound bool) {
	var (
		found  bool
		arg    int
		minVal float64
		minLbl int
	)

	// Find basic variable which limits change in entering variable.
	// If two choices result in the same change, the lower label must be preferred.
	for i := range dict.Basic {
		// Must have negative constraint coefficient.
		if dict.A[i][enter] >= -eps {
			continue
		}
		val := -dict.B[i] / dict.A[i][enter]
		lbl := dict.Basic[i]

		if found {
			if val > minVal {
				continue
			} else if val == minVal {
				if lbl > minLbl {
					continue
				}
			}
		}
		found = true
		arg = i
		minVal = val
		minLbl = lbl
	}
	return arg, !found
}

// NextBland returns the next pivot operation to perform according to Bland's rule.
// No pivot operation is possible if the dictionary is final or unbounded.
func NextBland(dict *Dict) Pivot {
	return NextBlandEps(dict, DefaultEps)
}

func NextBlandEps(dict *Dict, eps float64) Pivot {
	enter, final := toEnterBland(dict, eps)
	if final {
		return Pivot{Final: true}
	}
	leave, unbound := toLeaveBland(dict, enter, eps)
	if unbound {
		return Pivot{Unbounded: true}
	}
	return Pivot{Enter: enter, Leave: leave}
}

// NextFeasBland returns the next pivot operation to perform according to Bland's rule
// for a feasibility problem.
// This treats the variable with label len(NonBasic)+len(Basic)-1 as a special variable
// which receives priority to leave the basic set.
func NextFeasBland(dict *Dict) Pivot {
	return NextFeasBlandEps(dict, DefaultEps)
}

func NextFeasBlandEps(dict *Dict, eps float64) Pivot {
	m, n := len(dict.Basic), len(dict.NonBasic)
	// First check if there is a variable with label m+n-1 in the basic set.
	zero, found := find(m+n-1, dict.Basic)
	if found {
		// If there is and it can leave the basic set, make this pivot.
		enter, canLeave := toEnterFeasBland(dict, zero, eps)
		if canLeave {
			return Pivot{Enter: enter, Leave: zero}
		}
	}
	return NextBlandEps(dict, eps)
}

// Returns the non-basic variable to enter if the given variable were to leave.
// There may not exist such a non-basic variable.
func toEnterFeasBland(dict *Dict, leave int, eps float64) (enter int, found bool) {
	// Find min-label non-basic variable
	// which would choose the given basic variable to pivot with.
	var (
		arg int
		min int
	)

	for j := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[j] <= eps {
			continue
		}

		// Find leaving variable.
		// Element with minimum label chosen preferentially.
		i, unbnd := toLeaveBland(dict, j, eps)
		if unbnd {
			continue
		}

		if i == leave {
			lbl := dict.NonBasic[j]
			if !found || lbl < min {
				found = true
				arg = j
				min = dict.NonBasic[j]
			}
		}
	}
	return arg, found
}
