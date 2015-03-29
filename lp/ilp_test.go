package lp_test

import (
	"fmt"
	"log"

	"github.com/jvlmdr/golp/lp"
)

func Example() {
	const intTol = 1e-9

	dict := lp.NewDict(3, 2)
	dict.NonBasic = []int{0, 1}
	dict.Basic = []int{2, 3, 4}
	// max_{x, y >= 0} x + 2 y
	dict.C = []float64{1, 2}
	// subject to
	// -x + y <= 1,     x -  y +  1 >= 0
	// 3x + 2y <= 12, -3x - 2y + 12 >= 0
	// 2x + 3y <= 12, -2x - 3y + 12 >= 0
	dict.A[0], dict.B[0] = []float64{1, -1}, 1
	dict.A[1], dict.B[1] = []float64{-3, -2}, 12
	dict.A[2], dict.B[2] = []float64{-2, -3}, 12

	if !dict.Feas() {
		log.Println("init. dict. not feasible: solve feas. problem")
		// Solve the feasibility problem.
		var infeas bool
		dict, infeas = lp.SolveFeas(dict)
		if infeas {
			// Continuous relaxation is infeasible,
			// therefore integer problem is infeasible.
			log.Println("infeasible")
			return
		}
	} else {
		log.Println("init. dict. feasible")
	}

	log.Println("solve")
	var unbnd bool
	dict, unbnd = lp.Solve(dict)
	if unbnd {
		// Relaxation became unbounded.
		log.Println("primal is unbounded")
		return
	}

	for !dict.IsIntEps(intTol) {
		log.Printf("primal objective coeffs: %.4g\n", dict.C)

		// Check if integer solution.
		log.Println("add cutting-plane constraints")
		dict = lp.CutPlane(dict)
		log.Println("feasible?", dict.Feas())

		log.Println("switch to dual")
		dict = dict.Dual()
		log.Println("feasible?", dict.Feas())

		log.Println("solve")
		var unbnd bool
		dict, unbnd = lp.Solve(dict)
		if unbnd {
			// Dual of relaxation became unbounded.
			log.Println("dual is unbounded (primal infeasible)")
			return
		}

		log.Println("feasible?", dict.Feas())

		log.Println("switch to primal")
		dict = dict.Dual()
		log.Println("objective:", dict.Obj())
	}
	fmt.Printf("%.6g at %.6g\n", dict.Obj(), dict.Soln()[:2])
	// Output:
	// 6 at [2 2]
}
