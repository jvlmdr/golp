package main

import (
	"github.com/jackvalmadre/golp/lp"

	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	var (
		refFile string
		outFile string
	)
	flag.StringVar(&refFile, "ref", "", "File containing desired solution")
	flag.StringVar(&outFile, "out", "", "File to contain output")

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	dictFile := flag.Arg(0)

	// Load input dictionary.
	dict := new(lp.Dict)
	err := Load(lp.ReadDictColoradoFrom, dictFile, dict)
	if err != nil {
		log.Fatal(err)
	}

	const epsInt = 1e-3
	lp.DefaultEps = 1e-6

	result := solve(dict, epsInt)

	if outFile != "" {
		// Save dictionary out.
		if err := Save(WriteSolutionTo, outFile, result); err != nil {
			log.Fatalln("could not save solution")
		}
		log.Print("done")
	}

	if refFile != "" {
		// Load reference output.
		var ref Solution
		if err := Load(ReadSolutionFrom, refFile, &ref); err != nil {
			log.Fatalln("could not load reference:", err)
		}

		// Compare.
		if err := check(result, ref, epsInt); err != nil {
			log.Fatalln("fail:", err)
		}
		log.Print("pass")
	}
}

func solve(dict *lp.Dict, epsInt float64) Solution {
	if !dict.Feas() {
		log.Println("init. dict. not feasible: solve feas. problem")
		// Solve the feasibility problem.
		var infeas bool
		dict, infeas = lp.SolveFeas(dict)
		if infeas {
			// Continuous relaxation is infeasible,
			// therefore integer problem is infeasible.
			fmt.Println("infeasible")
			return Solution{Infeas: true}
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
		return Solution{Unbnd: true}
	}

	for !dict.IsIntEps(epsInt) {
		fmt.Printf("primal objective coeffs: %.4g\n", dict.C)

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
			return Solution{Infeas: true}
		}

		log.Println("feasible?", dict.Feas())

		log.Println("switch to primal")
		dict = dict.Dual()
		log.Println("objective:", dict.Obj())

		//	// Find a feasible solution.
		//	var infeas bool
		//	log.Println("solve feasibility problem")
		//	dict, infeas = lp.SolveFeas(dict)
		//	if infeas {
		//		// Continuous relaxation is infeasible,
		//		// therefore integer problem is infeasible.
		//		log.Println("primal is infeasible")
		//		return Solution{Infeas: true}
		//	}

		//	var unbnd bool
		//	log.Println("solve problem")
		//	dict, unbnd = lp.Solve(dict)
		//	if unbnd {
		//		// Relaxation became unbounded.
		//		log.Println("primal is unbounded")
		//		return Solution{Unbnd: true}
		//	}
	}

	return Solution{Obj: dict.Obj()}
}

func check(result, ref Solution, eps float64) error {
	if result.Infeas || ref.Infeas {
		if result.Infeas == ref.Infeas {
			// Success.
			return nil
		}
		return fmt.Errorf("infeasible: got %v, want %v", result.Infeas, ref.Infeas)
	}

	if result.Unbnd || ref.Unbnd {
		if result.Unbnd == ref.Unbnd {
			// Success.
			return nil
		}
		return fmt.Errorf("unbounded: got %v, want %v", result.Unbnd, ref.Unbnd)
	}

	// Check objective value.
	if math.Abs(result.Obj-ref.Obj) >= eps {
		return fmt.Errorf("objective: got %g, want %g", result.Obj, ref.Obj)
	}
	return nil
}
