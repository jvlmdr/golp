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

	result := solve(dict)

	if refFile != "" {
		// Load reference output.
		var ref Solution
		if err := Load(ReadSolutionFrom, refFile, &ref); err != nil {
			log.Fatalln("could not load reference:", err)
		}

		// Compare.
		if err := check(result, ref); err != nil {
			log.Fatalln("fail:", err)
		}
		log.Print("pass")
	}

	//	if outFile != "" {
	//		// Save dictionary out.
	//		if err := Save(WriteSolutionTo, outFile, result); err != nil {
	//			log.Fatalln("could not save solution")
	//		}
	//		log.Print("done")
	//	}
}

func solve(dict *lp.Dict) Solution {
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

	var dual bool
	for {
		// Solve.
		log.Println("solve")
		var unbnd bool
		dict, unbnd = lp.Solve(dict)
		if unbnd {
			// Relaxation became unbounded.
			if dual {
				log.Println("dual is unbounded (primal infeasible)")
				return Solution{Infeas: true}
			} else {
				log.Println("primal is unbounded")
				return Solution{Unbnd: true}
			}
		}

		// If the problem we just solved was the dual,
		// revert to the primal.
		if dual {
			dict = dict.Dual()
			dual = !dual
		}

		dict.Fprint(os.Stdout)
		fmt.Println()
		if dict.IsInt() {
			break
		}

		// Check if integer solution.
		log.Println("add cutting-plane constraints")
		dict = lp.CutPlane(dict)
		log.Println("feasible?", dict.Feas())

		log.Println("switch to dual of dictionary")
		dict = dict.Dual()
		dual = !dual
		log.Println("feasible?", dict.Feas())
	}

	// Dual of the dual is the primal.
	if dual {
		dict = dict.Dual()
		dual = !dual
	}

	dict.Fprint(os.Stdout)
	fmt.Println()

	return Solution{Obj: dict.Obj()}
}

func check(result, ref Solution) error {
	const eps = 1e-6

	// Check objective value.
	if math.Abs(result.Obj-ref.Obj) >= eps {
		return fmt.Errorf("objective: got %g, want %g", result.Obj, ref.Obj)
	}
	return nil
}
