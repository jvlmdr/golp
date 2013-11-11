package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Solution struct {
	Obj    float64
	Infeas bool
	Unbnd  bool
}

func ReadSolutionFrom(r io.Reader) (Solution, error) {
	scanner := bufio.NewScanner(r)
	var line string

	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	line = strings.TrimSpace(scanner.Text())

	if line == "infeasible" {
		return Solution{Infeas: true}, nil
	}
	if line == "unbounded" {
		return Solution{Unbnd: true}, nil
	}

	// Read a single number.
	obj, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return Solution{}, err
	}

	return Solution{Obj: obj}, nil
}

func WriteSolutionTo(w io.Writer, soln Solution) error {
	if soln.Infeas {
		_, err := fmt.Fprintln(w, "infeasible")
		return err
	}
	if soln.Unbnd {
		_, err := fmt.Fprintln(w, "unbounded")
		return err
	}

	if _, err := fmt.Fprintln(w, soln.Obj); err != nil {
		return err
	}
	return nil
}
