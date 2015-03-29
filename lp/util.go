package lp

import "math"

// Minimum of two integers.
func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// Maximum of two integers.
func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// Returns x - floor(x).
func mod1(x float64) float64 {
	return x - math.Floor(x)
}

// Returns x - round(x).
func distInt(x float64) float64 {
	return math.Min(x-math.Floor(x), math.Ceil(x)-x)
}

func find(x int, labels []int) (idx int, found bool) {
	for i, lbl := range labels {
		if lbl == x {
			return i, true
		}
	}
	return 0, false
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
