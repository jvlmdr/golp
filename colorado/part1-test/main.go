package main

import (
	"github.com/jackvalmadre/golp/lp"

	"fmt"
	"log"
	"math"
	"os"
	"path"
)

func main() {
	const (
		eps = 1e-3
		dir = "files"
		n   = 10
	)

	var ok int

	for i := 1; i <= n; i++ {
		fmt.Printf("case %d:\n", i)

		// Load dictionary and solution.
		dict, err := loadDict(path.Join(dir, fmt.Sprintf("dict%d", i)))
		if err != nil {
			log.Fatal(err)
		}
		soln, err := loadSoln(path.Join(dir, fmt.Sprintf("dict%d.output", i)))
		if err != nil {
			log.Fatal(err)
		}

		// Do test.
		if err := test(dict, soln, eps); err != nil {
			log.Print(err)
		}
		ok++
		fmt.Println("OK!")
	}

	fmt.Println("---")
	fmt.Printf("%d/%d passed\n", ok, n)
}

func test(dict *lp.Dict, soln Solution, eps float64) error {
	// Get entering variable, check final.
	enter, final := dict.Enter()
	if final {
		return fmt.Errorf("wrong: final: got %v, want %v", final, false)
	}

	// Get leaving variable, check unbounded.
	leave, unbounded := dict.Leave(enter)
	if unbounded != soln.Unbounded {
		return fmt.Errorf("wrong: unbounded: got %v, want %v", unbounded, soln.Unbounded)
	}
	// Only keep checking if bounded.
	if soln.Unbounded {
		return nil
	}

	if dict.NonBasic[enter] != soln.Enter {
		return fmt.Errorf("wrong enter index: got %d, want %d", dict.NonBasic[enter], soln.Enter)
	}
	if dict.Basic[leave] != soln.Leave {
		return fmt.Errorf("wrong leave index: got %d, want %d", dict.Basic[leave], soln.Leave)
	}

	dict = dict.Pivot(enter, leave)
	if !approx(dict.D, soln.Objective, eps) {
		return fmt.Errorf("wrong objective: got %g, want %g", dict.D, soln.Objective)
	}
	return nil
}

func approx(got, want, eps float64) bool {
	return math.Abs(got-want) <= eps*math.Abs(want)
}

func loadSoln(fname string) (Solution, error) {
	file, err := os.Open(fname)
	if err != nil {
		return Solution{}, err
	}
	defer file.Close()
	return ReadSolutionFrom(file)
}

func loadDict(fname string) (*lp.Dict, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return lp.ReadDictFrom(file)
}
