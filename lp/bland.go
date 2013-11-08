package lp

// Finds index (not label) of next non-basic variable to enter
// under Bland's rule.
func ToEnterBland(dict *Dict) (enter int, final bool) {
	// Find variable with lowest index.
	var (
		found bool
		arg   int
		min   int
	)

	// Find lowest-index variable with positive objective coefficient.
	for i := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[i] <= 0 {
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
func ToLeaveBland(dict *Dict, enter int) (leave int, unbound bool) {
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
		if dict.A[i][enter] >= 0 {
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

func NextBland(dict *Dict) Pivot {
	enter, final := ToEnterBland(dict)
	if final {
		return Pivot{Final: true}
	}
	leave, unbound := ToLeaveBland(dict, enter)
	if unbound {
		return Pivot{Unbounded: true}
	}
	return Pivot{Enter: enter, Leave: leave}
}
