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

func solve(dict *lp.Dict) Solution {
	for k := 0; ; k++ {
		// Get entering variable, check final.
		enter, final := dict.Enter()
		if final {
			return Solution{Value: dict.D, Steps: k}
		}

		// Get leaving variable, check unbounded.
		leave, unbounded := dict.Leave(enter)
		if unbounded {
			return Solution{Unbounded: true}
		}

		// Make the pivot.
		dict = dict.Pivot(enter, leave)
	}
}

func test(dict *lp.Dict, soln Solution, eps float64) error {
	result := solve(dict)

	if result.Unbounded != soln.Unbounded {
		return fmt.Errorf("wrong: unbounded: got %v, want %v", result.Unbounded, soln.Unbounded)
	}
	if soln.Unbounded {
		return nil
	}

	// Check objective value.
	if !approx(result.Value, soln.Value, eps) {
		return fmt.Errorf("wrong objective: got %g, want %g", result.Value, soln.Value)
	}

	// Check number of steps.
	if result.Steps != soln.Steps {
		return fmt.Errorf("wrong number of steps: got %d, want %d", result.Steps, soln.Steps)
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
