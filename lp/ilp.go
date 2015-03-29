package lp

import (
	"fmt"
	"log"
	"math"
)

// SolveInt solves the linear program with the constraint
// that all variables take integer values.
func SolveInt(dict *Dict) (final *Dict, err error) {
	return SolveIntEps(dict, DefaultEps)
}

func SolveIntEps(dict *Dict, eps float64) (final *Dict, err error) {
	if !dict.Feas() {
		// Solve the feasibility problem.
		var infeas bool
		dict, infeas = SolveFeasEps(dict, eps)
		if infeas {
			// Continuous relaxation is infeasible,
			// therefore integer problem is infeasible.
			return nil, fmt.Errorf("real problem is infeasible")
		}
	}

	// Solve problem without integer constraints.
	var unbnd bool
	dict, unbnd = SolveEps(dict, eps)
	if unbnd {
		// Relaxation became unbounded.
		return nil, fmt.Errorf("unbounded in primal")
	}

	for !dict.IsIntEps(eps) {
		log.Println("add cutting-plane constraints")
		dict = CutPlaneEps(dict, eps)
		log.Println("feasible?", dict.Feas())
		log.Println("switch to dual")
		dict = dict.Dual()
		log.Println("feasible?", dict.Feas())
		var unbnd bool
		dict, unbnd = SolveEps(dict, eps)
		if unbnd {
			// Dual of relaxation became unbounded.
			log.Println("unbounded in dual, infeasible in primal")
			return
		}
		log.Println("feasible?", dict.Feas())
		log.Println("switch to primal")
		dict = dict.Dual()
		log.Println("objective:", dict.Obj())
	}
	return dict, nil
}

// IsInt returns true if the dictionary is associated with an integer solution.
func (dict *Dict) IsInt() bool {
	return dict.IsIntEps(DefaultEps)
}

func (dict *Dict) IsIntEps(eps float64) bool {
	for i := range dict.Basic {
		// Distance from nearest int.
		dist := math.Abs(distInt(dict.B[i]))
		if dist > eps {
			// Basic variable is not (even roughly) an integer.
			return false
		}
	}
	return true
}

// CutPlane returns a new dictionary with cutting-plane constraints
// for non-integer expressions in the basic variables and objective.
func CutPlane(orig *Dict) *Dict {
	return CutPlaneEps(orig, DefaultEps)
}

func CutPlaneEps(orig *Dict, eps float64) *Dict {
	var A [][]float64
	var B []float64
	m, n := len(orig.Basic), len(orig.NonBasic)

	for i := range orig.Basic {
		// Distance from nearest int.
		dist := math.Abs(distInt(orig.B[i]))
		if dist <= eps {
			// Basic variable is (roughly) an integer.
			continue
		}
		// Not (even roughly) an integer.
		// Add cutting plane constraint.
		a := make([]float64, n)
		for j := 0; j < n; j++ {
			a[j] = mod1(-orig.A[i][j])
		}
		b := -mod1(orig.B[i])
		A = append(A, a)
		B = append(B, b)
	}

	// Do same for objective.

	// Distance from nearest int.
	dist := math.Abs(distInt(orig.D))
	if dist <= eps {
		// Basic variable is (roughly) an integer.
	} else {
		// Not (even roughly) an integer.
		// Add cutting plane constraint.
		a := make([]float64, n)
		for j := 0; j < n; j++ {
			a[j] = mod1(-orig.C[j])
		}
		b := -mod1(orig.D)
		A = append(A, a)
		B = append(B, b)
	}

	// Copy dictionary.
	dict := NewDict(m+len(A), n)
	copy(dict.Basic, orig.Basic)
	copy(dict.NonBasic, orig.NonBasic)
	for i := range orig.A {
		copy(dict.A[i], orig.A[i])
	}
	copy(dict.B, orig.B)
	copy(dict.C, orig.C)
	dict.D = orig.D

	// Add new rows and slack variables.
	for i := range A {
		dict.Basic[m+i] = m + n + i
		dict.A[m+i] = A[i]
		dict.B[m+i] = B[i]
	}
	return dict
}
