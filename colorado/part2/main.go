package main

import (
	"github.com/jackvalmadre/golp/lp"

	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
)

func main() {
	const (
		dir = "files"
		n   = 5
	)

	for i := 1; i <= n; i++ {
		fmt.Printf("case %d:\n", i)

		// Load dictionary and solution.
		dictFile := path.Join(dir, fmt.Sprintf("part%d.dict", i))
		dict, err := loadDict(dictFile)
		if err != nil {
			log.Fatal(err)
		}

		soln := solve(dict)

		solnFile := path.Join(dir, fmt.Sprintf("part%d.output", i))
		if err := saveSolutionTo(solnFile, soln); err != nil {
			log.Fatal(err)
		}
	}
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

func saveSolutionTo(fname string, soln Solution) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeSolutionTo(file, soln)
}

func writeSolutionTo(w io.Writer, soln Solution) error {
	if soln.Unbounded {
		if _, err := fmt.Fprintln(w, "UNBOUNDED"); err != nil {
			return err
		}
		return nil
	}

	if _, err := fmt.Fprintln(w, soln.Value); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, soln.Steps); err != nil {
		return err
	}
	return nil
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
