package lp

var DefaultEps = 1e-9

// Dict is a "dictionary" describing a linear program.
// The simplex algorithm moves from dictionary to dictionary.
// The problem is
//	max {D + sum_i C[i] x[NonBasic[i]]}
//	s.t. x[NonBasic[i]] >= 0
//	     x[Basic[i]] = B[j] + sum_j A[i][j] x[NonBasic[i]] >= 0.
// The variables are partitioned into basic and non-basic sets.
// The objective and the basic variables are affine functions
// of the non-basic variables.
// All variables are constrained to be greater than zero.
// The number of basic variables is the number of affine inequalities.
//
// Each dictionary is associated with the "solution"
// where all of its non-basic variables are zero.
// Note that this is not necessarily the maximizer or even feasible.
//
// A dictionary is feasible if all B[i] >= 0
// since this implies that if all non-basic variables are zero,
// then all non-basic variables are non-negative.
type Dict struct {
	Basic    []int
	NonBasic []int
	// basic = A nonbasic + b
	A [][]float64
	B []float64
	// objective = c' nonbasic + d
	C []float64
	D float64
}

// NewDict creates a dictionary with m basic and n non-basic variables.
func NewDict(m, n int) *Dict {
	d := new(Dict)
	d.Basic = make([]int, m)
	d.NonBasic = make([]int, n)
	d.A = make([][]float64, m)
	for i := range d.A {
		d.A[i] = make([]float64, n)
	}
	d.B = make([]float64, m)
	d.C = make([]float64, n)
	return d
}

// Soln returns the solution associated with the dictionary.
// This is the value of all variables when the non-basic variables are zero.
func (dict *Dict) Soln() []float64 {
	p := len(dict.Basic) + len(dict.NonBasic)
	x := make([]float64, p)
	// Non-basic are all zero (default value).
	// Therefore basic = b.
	for i, j := range dict.Basic {
		x[j] = dict.B[i]
	}
	return x
}

// Obj returns the objective value associated with the solution of this dictionary.
func (dict *Dict) Obj() float64 {
	return dict.D
}

// Feas returns true if the (solution associated with the) dictionary is feasible.
func (dict *Dict) Feas() bool {
	return dict.FeasEps(DefaultEps)
}

func (dict *Dict) FeasEps(eps float64) bool {
	// Infeasible if any of the basic variables are less than zero.
	for _, bi := range dict.B {
		// Let them be very small and negative.
		if bi < -eps {
			return false
		}
	}
	return true
}

// Pivot swaps Basic[leave] and NonBasic[enter].
func (src *Dict) Pivot(enter, leave int) *Dict {
	m := len(src.Basic)
	n := len(src.NonBasic)
	dst := NewDict(m, n)

	// Copy the variable indices.
	copy(dst.Basic, src.Basic)
	copy(dst.NonBasic, src.NonBasic)
	// Swap variables to enter and leave.
	dst.Basic[leave], dst.NonBasic[enter] = src.NonBasic[enter], src.Basic[leave]

	//	log.Printf(
	//		"pivot: enter %d, leave %d, obj coeff %.4g, basic coeff %.4g",
	//		src.NonBasic[enter], src.Basic[leave],
	//		src.C[enter], src.A[leave][enter],
	//	)

	// Update row of basic variable.
	dst.B[leave] = -src.B[leave] / src.A[leave][enter]
	for j := range dst.A[leave] {
		if j == enter {
			dst.A[leave][j] = 1 / src.A[leave][enter]
		} else {
			dst.A[leave][j] = -src.A[leave][j] / src.A[leave][enter]
		}
	}

	// Update column of non-basic variable.
	for i := range dst.A {
		if i == leave {
			continue
		}
		dst.B[i] = src.B[i] + src.A[i][enter]*dst.B[leave]
		for j := range dst.A[i] {
			if j == enter {
				dst.A[i][j] = src.A[i][enter] / src.A[leave][j]
			} else {
				dst.A[i][j] = src.A[i][j] + src.A[i][enter]*dst.A[leave][j]
			}
		}
	}

	// Update objective row.
	dst.D = src.D + src.C[enter]*dst.B[leave]
	for j := range dst.C {
		if j == enter {
			dst.C[j] = src.C[enter] / src.A[leave][j]
		} else {
			dst.C[j] = src.C[j] + src.C[enter]*dst.A[leave][j]
		}
	}
	return dst
}

// ToFeasDict creates a dictionary describing the feasibility problem.
// Adds a variable to the basic set, then pivots it into the non-basic set.
// Assumes that Basic and NonBasic are indices from 0 to len(NonBasic)+len(Basic)-1.
func ToFeasDict(infeas *Dict) *Dict {
	m, n := len(infeas.Basic), len(infeas.NonBasic)
	// Add a new non-basic variable.
	dict := NewDict(m, n+1)

	// Copy constraints.
	copy(dict.Basic, infeas.Basic)
	copy(dict.NonBasic, infeas.NonBasic)
	dict.NonBasic[n] = m + n
	for i := 0; i < m; i++ {
		copy(dict.A[i], infeas.A[i])
	}
	copy(dict.B, infeas.B)

	// Add new non-basic variable to all rows.
	for i := 0; i < m; i++ {
		dict.A[i][n] = 1
	}
	// Goal is to minimize this new variable.
	dict.C[n] = -1

	leave := findMin(dict.B)
	return dict.Pivot(n, leave)
}

// FromFeasDict returns to original problem.
// Removes the variable with label len(NonBasic)+len(Basic)-1,
// which must be in the non-basic set.
func FromFeasDict(feas *Dict, orig *Dict) *Dict {
	m, n := len(feas.Basic), len(feas.NonBasic)-1

	// Remove extra variable.
	// Must be non-basic.
	extra, found := find(m+n-1, feas.NonBasic)
	if !found {
		panic("extra variable not in non-basic set")
	}

	dict := NewDict(m, n)
	// Copy constraints.
	copy(dict.Basic, feas.Basic)
	copy(dict.NonBasic[:extra], feas.NonBasic[:extra])
	copy(dict.NonBasic[extra:], feas.NonBasic[extra+1:])
	for i := 0; i < m; i++ {
		copy(dict.A[i][:extra], feas.A[i][:extra])
		copy(dict.A[i][extra:], feas.A[i][extra+1:])
	}
	copy(dict.B, feas.B)

	// Re-express original objective in terms of current basic set.
	// This could be done succinctly with matrix operations?
	dict.D = orig.D
	for u, lbl1 := range orig.NonBasic {
		c := orig.C[u]
		for j, lbl2 := range dict.NonBasic {
			if lbl1 != lbl2 {
				continue
			}
			// Simply transfer coefficient.
			dict.C[j] += c
		}
		for i, lbl2 := range dict.Basic {
			if lbl1 != lbl2 {
				continue
			}
			// Transfer coefficients for basic variable.
			dict.D += c * dict.B[i]
			for j := range dict.NonBasic {
				dict.C[j] += c * dict.A[i][j]
			}
		}
	}
	return dict
}
