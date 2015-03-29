package lp_test

import (
	"fmt"

	"github.com/jvlmdr/golp/lp"
)

func ExampleSolve() {
	dict := new(lp.Dict)
	dict.NonBasic = []int{0, 1}
	dict.Basic = []int{2, 3, 4}
	// max_{x, y >= 0} x + 2 y
	dict.C = []float64{1, 2}
	// subject to
	// -x + y <= 1,     x -  y +  1 >= 0
	// 3x + 2y <= 12, -3x - 2y + 12 >= 0
	// 2x + 3y <= 12, -2x - 3y + 12 >= 0
	dict.A = make([][]float64, 3)
	dict.B = make([]float64, 3)
	dict.A[0], dict.B[0] = []float64{1, -1}, 1
	dict.A[1], dict.B[1] = []float64{-3, -2}, 12
	dict.A[2], dict.B[2] = []float64{-2, -3}, 12

	dict, unbnd := lp.Solve(dict)
	if unbnd {
		fmt.Print("unbounded")
		return
	}
	fmt.Printf("%.6g at %.6g\n", dict.Obj(), dict.Soln()[:2])
	// Output:
	// 7.4 at [1.8 2.8]
}
