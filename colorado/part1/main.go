package main

import (
	"flag"
	"fmt"
	"log"
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

	// Make the pivot.
	var soln Solution
	piv := lp.NextBland(dict)
	if piv.Unbounded {
		// The initial dictionary was unbounded.
		soln.Unbounded = true
	} else {
		if piv.Final {
			// The initial dictionary was final.
			soln.Final = true
		} else {
			// Possible to make a pivot.
			soln.Enter = dict.NonBasic[piv.Enter]
			soln.Leave = dict.Basic[piv.Leave]
			soln.Dict = dict.Pivot(piv.Enter, piv.Leave)
		}
	}

	if refFile != "" {
		// Load reference output.
		var ref Summary
		if err := Load(ReadSummaryFrom, refFile, &ref); err != nil {
			log.Fatalln("could not load reference summary:", err)
		}

		// Compare.
		ok := check(soln, ref)
		if !ok {
			log.Fatal("fail")
		}
		log.Print("pass")
	}

	if outFile != "" {
		// Save dictionary out.
		if err := Save(WriteSolutionTo, outFile, soln); err != nil {
			log.Fatalln("could not save solution")
		}
	}
}

func check(soln Solution, ref Summary) bool {
	// First check unbounded.
	if ref.Unbounded || soln.Unbounded {
		if ref.Unbounded != soln.Unbounded {
			log.Printf("unbounded: want %v, got %v", ref.Unbounded, soln.Unbounded)
			return false
		} else {
			return true
		}
	}

	// Otherwise check pivot indices.
	if ref.Enter != soln.Enter || ref.Leave != soln.Leave {
		log.Printf(
			"pivot: want %d enter and %d leave, got %d enter and %d leave",
			ref.Enter, ref.Leave, soln.Enter, soln.Leave,
		)
		return false
	}

	refObj := fmt.Sprintf("%.2f", ref.Objective)
	solnObj := fmt.Sprintf("%.2f", soln.Dict.Obj())
	if refObj != solnObj {
		log.Printf("objective: want %s, got %s", refObj, solnObj)
		return false
	}
	return true
}
