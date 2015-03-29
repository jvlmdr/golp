package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jvlmdr/golp/lp"
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

	// Pivot all the way.
	var iter int
	var unbounded bool
	for {
		piv := lp.NextBland(dict)
		if piv.Unbounded {
			unbounded = true
			fmt.Println("unbounded")
			break
		}
		if piv.Final {
			fmt.Println("final")
			break
		}
		dict = dict.Pivot(piv.Enter, piv.Leave)
		iter++
		fmt.Printf("%4d  f:%10.3e\n", iter, dict.Obj())
	}

	var result Solution
	if unbounded {
		result = Solution{Unbounded: true}
	} else {
		result = Solution{Value: dict.Obj(), Steps: iter}
	}

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
	}
}

func check(result, ref Solution) error {
	const eps = 0.1

	if result.Unbounded || ref.Unbounded {
		if result.Unbounded != ref.Unbounded {
			return fmt.Errorf("unbounded: got %v, want %v", result.Unbounded, ref.Unbounded)
		} else {
			// Both are true.
			return nil
		}
	}

	// Check objective value.
	if math.Abs(result.Value-ref.Value) >= eps {
		return fmt.Errorf("objective: got %g, want %g", result.Value, ref.Value)
	}
	// Check number of steps.
	if result.Steps != ref.Steps {
		return fmt.Errorf("number of steps: got %d, want %d", result.Steps, ref.Steps)
	}
	return nil
}
