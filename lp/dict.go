package lp

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

// Returns the solution associated with the dictionary.
// The non-basic variables are zero.
func (dict *Dict) Soln() []float64 {
	n := len(dict.Basic) + len(dict.NonBasic)
	x := make([]float64, n)
	// Non-basic are all zero (default value).
	// Therefore basic = b.
	for _, idx := range dict.Basic {
		x[idx] = dict.B[idx]
	}
	return x
}

// Returns the objective value.
func (dict *Dict) Obj() float64 {
	return dict.D
}

// Is the (solution associated with the) dictionary feasible?
func (dict *Dict) Feas() bool {
	for _, bi := range dict.B {
		if bi < 0 {
			return false
		}
	}
	return true
}

// Swaps the non-basic variable to enter and the basic variable to leave.
func (src *Dict) Pivot(enter, leave int) *Dict {
	m := len(src.Basic)
	n := len(src.NonBasic)
	dst := NewDict(m, n)

	// Copy the variable indices.
	copy(dst.Basic, src.Basic)
	copy(dst.NonBasic, src.NonBasic)
	// Swap variables to enter and leave.
	dst.Basic[leave] = src.NonBasic[enter]
	dst.NonBasic[enter] = src.Basic[leave]

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

// Creates a dictionary describing the feasibility problem.
func FeasDict(infeas *Dict) *Dict {
	m, n := len(infeas.Basic), len(infeas.NonBasic)
	// Add a new non-basic variable.
	dict := NewDict(m, n+1)

	// Copy constraints.
	copy(dict.Basic, infeas.Basic)
	copy(dict.NonBasic, infeas.NonBasic)
	for i := 0; i < m; i++ {
		copy(dict.A[i], infeas.A[i])
	}
	copy(dict.B, infeas.B)

	// Add new non-basic variable (with label 0) to all rows.
	for i := 0; i < m; i++ {
		dict.A[i][n] = 1
	}
	// Goal is to minimize this new variable.
	dict.C[n] = -1

	leave := findMin(dict.B)
	return dict.Pivot(n, leave)
}

func findMin(vals []float64) int {
	var arg int
	for i, v := range vals {
		if v < vals[arg] {
			arg = i
		}
	}
	return arg
}
