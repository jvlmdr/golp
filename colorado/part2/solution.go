package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Solution struct {
	Value     float64
	Steps     int
	Unbounded bool
}

func ReadSolutionFrom(r io.Reader) (Solution, error) {
	scanner := bufio.NewScanner(r)
	var line string

	// Read first line.
	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	// First line contains UNBOUNDED or a number.
	line = strings.TrimSpace(scanner.Text())
	if line == "UNBOUNDED" {
		return Solution{Unbounded: true}, nil
	}
	value, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return Solution{}, err
	}

	// Read second line.
	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	// Second line contains an integer.
	line = strings.TrimSpace(scanner.Text())
	steps, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return Solution{}, err
	}

	soln := Solution{Value: value, Steps: int(steps)}
	return soln, nil
}

func WriteSolutionTo(w io.Writer, soln Solution) error {
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
