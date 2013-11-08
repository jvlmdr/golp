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

	dict.Fprint(os.Stdout)
	fmt.Println("feas?", dict.Feas())
	fmt.Println()

	dict = lp.FeasDict(dict)
	dict.Fprint(os.Stdout)
	fmt.Println("feas?", dict.Feas())
	fmt.Println()

	// Pivot all the way.
	var iter int
	for {
		piv := lp.NextFeasBland(dict)
		if piv.Unbounded {
			panic("unbounded")
		}
		if piv.Final {
			fmt.Println("final")
			break
		}
		dict = dict.Pivot(piv.Enter, piv.Leave)
		iter++
		fmt.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}

	result := Solution{dict.Obj()}

	if refFile != "" {
		// Load reference output.
		var ref Solution
		if err := Load(ReadSolutionFrom, refFile, &ref); err != nil {
			log.Fatalln("could not load reference:", err)
		}

		// Compare.
		if err := check(result, ref); err != nil {
			log.Println("fail:", err)
		}
		log.Print("pass")
	}

	if outFile != "" {
		// Save dictionary out.
		if err := Save(WriteSolutionTo, outFile, result); err != nil {
			log.Fatalln("could not save solution")
		}
		log.Print("done")
	}
}

func check(result, ref Solution) error {
	const eps = 1e-6

	// Check objective value.
	if math.Abs(result.Obj-ref.Obj) >= eps {
		return fmt.Errorf("objective: got %g, want %g", result.Obj, ref.Obj)
	}
	return nil
}
