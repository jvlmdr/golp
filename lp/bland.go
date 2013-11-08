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
func toLeaveBland(dict *Dict, enter int) (leave int, unbound bool) {
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
	leave, unbound := toLeaveBland(dict, enter)
	if unbound {
		return Pivot{Unbounded: true}
	}
	return Pivot{Enter: enter, Leave: leave}
}

func NextFeasBland(dict *Dict) Pivot {
	// First check if there is a variable with label 0 in the basic set.
	zero, foundZero := findZero(dict.Basic)
	if foundZero {
		// If there is and it can leave the basic set, make this pivot.
		enter, canLeave := toEnterFeasBland(dict, zero)
		if canLeave {
			return Pivot{Enter: enter, Leave: zero}
		}
	}

	return NextBland(dict)
}

func findZero(labels []int) (idx int, found bool) {
	for i, lbl := range labels {
		if lbl == 0 {
			return i, true
		}
	}
	return 0, false
}

func toEnterFeasBland(dict *Dict, leave int) (enter int, found bool) {
	// Find variable with lowest index which can enter.

	var (
		arg int
		min int
	)

	// Find lowest-index variable with positive objective coefficient.
	for j := range dict.NonBasic {
		// Must have positive objective coefficient.
		if dict.C[j] <= 0 {
			continue
		}

		// Find leaving variable.
		// Element with minimum label chosen preferentially.
		i, unbnd := toLeaveBland(dict, j)
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
