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
	Unbnd  bool
	Infeas bool
}

func ReadSolutionFrom(r io.Reader) (Solution, error) {
	scanner := bufio.NewScanner(r)
	var line string

	// Read a single number.
	if err := readLine(scanner); err != nil {
		return Solution{}, err
	}
	line = strings.TrimSpace(scanner.Text())
	obj, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return Solution{}, err
	}

	return Solution{Obj: obj}, nil
}

func WriteSolutionTo(w io.Writer, soln Solution) error {
	if _, err := fmt.Fprintln(w, soln.Obj); err != nil {
		return err
	}
	return nil
}
