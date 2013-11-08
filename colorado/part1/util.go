package main

import "math"

func eqEps(x, ref, eps float64) bool {
	rel := math.Abs(x-ref) / math.Abs(ref)
	return rel < eps
}
