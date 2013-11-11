package lp

import (
	"math"
)

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

// Incorporates cutting-plane constraints
// for all non-integer variables and possibly the objective.
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
		dict.Basic[m+i] = m+n+1+i
		dict.A[m+i] = A[i]
		dict.B[m+i] = B[i]
	}
	return dict
}
